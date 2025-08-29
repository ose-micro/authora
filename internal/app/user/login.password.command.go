package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/domain/role"
	"github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/common"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/core/utils"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (u loginCommandHandler) Handle(ctx context.Context, command user.LoginCommand) (*user.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.login.command.handler", trace.WithAttributes(
		attribute.String("operation", "login"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command payload
	if err := command.Validate(); err != nil {
		err := ose_error.New(ose_error.ErrInvalidInput, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := u.repo.User.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "email",
						Op:    dto.OpEq,
						Value: command.Email,
					},
				},
			},
		},
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to login record",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	if utils.CheckPasswordHash(record.Password(), command.Password) {
		err := ose_error.New(ose_error.ErrUnauthorized, "password does not match")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to login record",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	// Fetch all roles for user
	assigns, err := u.repo.Assignment.Read(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "user",
						Op:    dto.OpEq,
						Value: record.ID(),
					},
				},
			},
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to login record",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	var assignmentTenants []role.Public
	if err := common.JsonToAny(assigns["one"], assignmentTenants); err != nil {
		return nil, err
	}

	fmt.Println(assignmentTenants)

	// save user to write store
	err = u.repo.User.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to login record",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info("login process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "login"),
		zap.Any("payload", command),
	)

	return record, nil
}

func newLoginCommandHandler(bs domain.Domain, repo repository.Repository,
	log logger.Logger, tracer tracing.Tracer) cqrs.CommandHandle[user.LoginCommand, *user.Domain] {
	return &loginCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

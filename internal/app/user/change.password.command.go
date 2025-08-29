package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type changePasswordCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (u changePasswordCommandHandler) Handle(ctx context.Context, command user.ChangePasswordCommand) (*user.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.change_password.command.handler", trace.WithAttributes(
		attribute.String("operation", "change_password"),
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
			zap.String("operation", "change_password"),
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
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Id,
					},
				},
			},
		},
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to change_password record",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_password"),
			zap.Error(err),
		)

		return nil, err
	}

	if err := record.ChangePassword(command.Password, command.OldPassword); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to change_password record",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_password"),
			zap.Error(err),
		)

		return nil, err
	}

	// save user to write store
	err = u.repo.User.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to change_password record",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_password"),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info("change_password process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "change_password"),
		zap.Any("payload", command),
	)

	return record, nil
}

func newChangePasswordCommandHandler(bs domain.Domain, repo repository.Repository,
	log logger.Logger, tracer tracing.Tracer) cqrs.CommandHandle[user.ChangePasswordCommand, *user.Domain] {
	return &changePasswordCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

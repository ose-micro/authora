package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
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

type resetPasswordCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (u resetPasswordCommandHandler) Handle(ctx context.Context, command user.ResetPasswordCommand) (*user.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.change_password.command.handler", trace.WithAttributes(
		attribute.String("operation", "change_password"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command dto
	if err := command.Validate(); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
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

	if !record.Status().IsActive() {
		err := ose_error.New(ose_error.ErrUnauthorized, "user is not active", traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to change_password record",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_password"),
			zap.Error(err),
		)

		return nil, err
	}

	if err := record.ResetPassword(command.NewPassword); err != nil {
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
		zap.Any("dto", command),
	)

	return record, nil
}

func newResetPasswordCommandHandler(bs business.Domain, repo repository.Repository,
	log logger.Logger, tracer tracing.Tracer) cqrs.CommandHandle[user.ResetPasswordCommand, *user.Domain] {
	return &resetPasswordCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

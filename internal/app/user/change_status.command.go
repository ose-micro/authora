package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/user"
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

const ChangeStatusOperation = "change_status"

type changeStatusCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (c changeStatusCommandHandler) Handle(ctx context.Context, command user.StatusCommand) (*bool, error) {
	ctx, span := c.tracer.Start(ctx, "app.user.create.command.handler", trace.WithAttributes(
		attribute.String("operation", ChangeStatusOperation),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command dto
	if err := command.Validate(); err != nil {
		err := ose_error.New(ose_error.ErrInvalidInput, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", ChangeStatusOperation),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := c.repo.User.ReadOne(ctx, dto.Request{
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
		c.log.Error("user role not found",
			zap.String("trace_id", traceId),
			zap.String("operation", ChangeStatusOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if err := record.ChangeState(command.State); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to change state",
			zap.String("trace_id", traceId),
			zap.String("operation", ChangeStatusOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if err := c.repo.User.Update(ctx, *record); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create user",
			zap.String("trace_id", traceId),
			zap.String("operation", ChangeStatusOperation),
			zap.Error(err),
		)

		return nil, err
	}

	c.log.Info("change status process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", ChangeStatusOperation),
		zap.Any("dto", command),
	)

	result := true

	return &result, nil
}

func newChangeStausCommandHandler(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.CommandHandle[user.StatusCommand, *bool] {
	return &changeStatusCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

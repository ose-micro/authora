package assignment

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/domain/assignment"
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

type updateCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (u updateCommandHandler) Handle(ctx context.Context, command assignment.UpdateCommand) (*assignment.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.assignment.update.command.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
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
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := u.repo.Assignment.ReadOne(ctx, dto.Request{
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
		u.log.Error("failed to update record",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	record.UpdateRole(command.Role)

	// save assignment to write store
	err = u.repo.Assignment.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to update record",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info("update process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("payload", command),
	)

	return record, nil
}

func newUpdateCommandHandler(bs domain.Domain, repo repository.Repository,
	log logger.Logger, tracer tracing.Tracer) cqrs.CommandHandle[assignment.UpdateCommand, *assignment.Domain] {
	return &updateCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

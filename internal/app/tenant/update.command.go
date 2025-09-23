package tenant

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/tenant"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type updateCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (u updateCommandHandler) Handle(ctx context.Context, command tenant.UpdateCommand) (*tenant.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.tenant.update.command.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command dto
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := u.repo.Tenant.ReadOne(ctx, dto.Request{
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

	record.Update(tenant.Params{
		Name:     command.Name,
		Metadata: command.Metadata,
	})

	// save tenant to write store
	err = u.repo.Tenant.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("fail while saving record",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info("update process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", command),
	)

	return record, nil
}

func newUpdateCommandHandler(bs business.Domain, repo repository.Repository,
	log logger.Logger, tracer tracing.Tracer) cqrs.CommandHandle[tenant.UpdateCommand, *tenant.Domain] {
	return &updateCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

package tenant

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/tenant"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const CreateOperation = "create"

type createCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (c createCommandHandler) Handle(ctx context.Context, command tenant.CreateCommand) (*tenant.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.tenant.create.command.handler", trace.WithAttributes(
		attribute.String("operation", CreateOperation),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command dto
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if check, _ := c.repo.Tenant.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "name",
						Op:    dto.OpEq,
						Value: command.Name,
					},
				},
			},
		},
	}); check != nil {
		err := fmt.Errorf("tenant with id %s already exists", command.Name)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := c.bs.Tenant.New(tenant.Params{
		Name:     command.Name,
		Metadata: command.Metadata,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error(err.Error(),
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if err := c.repo.Tenant.Create(ctx, *record); err != nil {
		return nil, err
	}

	c.log.Info("create process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", CreateOperation),
		zap.Any("dto", command),
	)

	return record, nil
}

func newCreateCommandHandler(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.CommandHandle[tenant.CreateCommand, *tenant.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

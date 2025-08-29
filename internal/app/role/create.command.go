package role

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/domain/role"
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

const CreateOperation = "create"

type createCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (c createCommandHandler) Handle(ctx context.Context, command role.CreateCommand) (*role.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.role.create.command.handler", trace.WithAttributes(
		attribute.String("operation", CreateOperation),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		err := ose_error.New(ose_error.ErrInvalidInput, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if check, _ := c.repo.Role.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "name",
						Op:    dto.OpEq,
						Value: command.Name,
					},
					{
						Field: "tenant",
						Op:    dto.OpEq,
						Value: command.Tenant,
					},
				},
			},
		},
	}); check != nil {
		err := ose_error.New(ose_error.ErrAlreadyExists, "role already exists")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if _, err := c.repo.Tenant.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Tenant,
					},
				},
			},
		},
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to check if tenant exists",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := c.bs.Role.New(role.Params{
		Name:        command.Name,
		Tenant:      command.Tenant,
		Permissions: command.Permissions,
	})
	if err != nil {
		err = ose_error.New(ose_error.ErrInternal, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create role",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if err := c.repo.Role.Create(ctx, *record); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create role",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	c.log.Info("create process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", CreateOperation),
		zap.Any("payload", command),
	)

	return record, nil
}

func newCreateCommandHandler(bs domain.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.CommandHandle[role.CreateCommand, *role.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

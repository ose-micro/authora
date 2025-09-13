package assignment

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/assignment"
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
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (c createCommandHandler) Handle(ctx context.Context, command assignment.CreateCommand) (*assignment.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.assignment.create.command.handler", trace.WithAttributes(
		attribute.String("operation", CreateOperation),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command dto
	if err := command.Validate(); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
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
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if _, err := c.repo.Role.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Role,
					},
					{
						Field: "tenant",
						Op:    dto.OpEq,
						Value: command.Tenant,
					},
				},
			},
		},
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if check, _ := c.repo.Assignment.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "user",
						Op:    dto.OpEq,
						Value: command.User,
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
		err := ose_error.New(ose_error.ErrConflict, "assignment already exists for this user on this tenant", traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := c.bs.Assignment.New(assignment.Params{
		User:   command.User,
		Tenant: command.Tenant,
		Role:   command.Role,
	})
	if err != nil {
		err = ose_error.New(ose_error.ErrBadRequest, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create assignment",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if err := c.repo.Assignment.Create(ctx, *record); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create assignment",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

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
	tracer tracing.Tracer) cqrs.CommandHandle[assignment.CreateCommand, *assignment.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

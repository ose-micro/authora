package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/domain"
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
	bus    domain.Bus
}

// Handle implements cqrs.CommandHandle.
func (c createCommandHandler) Handle(ctx context.Context, command user.CreateCommand) (*user.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.user.create.command.handler", trace.WithAttributes(
		attribute.String("operation", CreateOperation),
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
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	role, err := c.repo.Role.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Role,
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
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if check, _ := c.repo.User.ReadOne(ctx, dto.Request{
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
	}); check != nil {
		err := ose_error.New(ose_error.ErrAlreadyExists, "user already exists with this email")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := c.bs.User.New(user.Params{
		Email:      command.Email,
		GivenNames: command.GivenNames,
		FamilyName: command.FamilyName,
		Password:   command.Password,
	})
	if err != nil {
		err = ose_error.New(ose_error.ErrInternal, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create user",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	if err := c.repo.User.Create(ctx, *record); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create user",
			zap.String("trace_id", traceId),
			zap.String("operation", CreateOperation),
			zap.Error(err),
		)

		return nil, err
	}

	err = c.publishEvent(*record.Public(), *role)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to publish event",
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

func (c createCommandHandler) publishEvent(payload user.Public, role role.Domain) error {
	err := c.bus.Publish(user.CreatedEvent, user.DefaultEvent{
		ID:         payload.Id,
		GivenNames: payload.GivenNames,
		FamilyName: payload.FamilyName,
		Email:      payload.Email,
		Password:   payload.Password,
		Role:       role.ID(),
		Tenant:     role.Tenant(),
		Metadata:   payload.Metadata,
		CreatedAt:  payload.CreatedAt,
	})
	if err != nil {
		return ose_error.New(ose_error.ErrInternal, err.Error())
	}

	return nil
}

func newCreateCommandHandler(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer, bus domain.Bus) cqrs.CommandHandle[user.CreateCommand, *user.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
		bus:    bus,
	}
}

package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	ose_jwt "github.com/ose-micro/jwt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type hasRoleCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     business.Domain
	jwt    ose_jwt.Manager
}

// Handle implements cqrs.CommandHandle.
func (u hasRoleCommandHandler) Handle(ctx context.Context, command user.HasRoleCommand) (*bool, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.has_role.command.handler", trace.WithAttributes(
		attribute.String("operation", "has_role"),
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
			zap.String("operation", "has_role"),
			zap.Error(err),
		)

		return nil, err
	}

	claims, err := u.jwt.ParseClaims(command.Token)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrUnauthorized, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_role"),
			zap.Error(err),
		)

		return nil, err
	}
	result := ose_jwt.HasTenantRole(*claims, command.Tenant, command.Role)

	return &result, nil
}

func newHasRoleCommandHandler(log logger.Logger, tracer tracing.Tracer,
	jwt ose_jwt.Manager) cqrs.CommandHandle[user.HasRoleCommand, *bool] {
	return &hasRoleCommandHandler{
		log:    log,
		tracer: tracer,
		jwt:    jwt,
	}
}

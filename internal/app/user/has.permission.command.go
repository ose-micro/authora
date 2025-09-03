package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/user"
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

type hasPermissionCommandHandler struct {
	log    logger.Logger
	tracer tracing.Tracer
	jwt    ose_jwt.Manager
}

// Handle implements cqrs.CommandHandle.
func (u hasPermissionCommandHandler) Handle(ctx context.Context, command user.HasPermissionCommand) (*bool, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.has_permission.command.handler", trace.WithAttributes(
		attribute.String("operation", "has_permission"),
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
			zap.String("operation", "has_permission"),
			zap.Error(err),
		)

		return nil, err
	}

	claims, err := u.jwt.ParseClaims(command.Token)
	if err != nil {
		err := ose_error.New(ose_error.ErrUnauthorized, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_permission"),
			zap.Error(err),
		)

		return nil, err
	}
	result := ose_jwt.HasTenantPermission(*claims, command.Tenant, *command.Permission)
	return &result, nil
}

func newHasPermissionCommandHandler(log logger.Logger, tracer tracing.Tracer,
	jwt ose_jwt.Manager) cqrs.CommandHandle[user.HasPermissionCommand, *bool] {
	return &hasPermissionCommandHandler{
		log:    log,
		tracer: tracer,
		jwt:    jwt,
	}
}

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

type parseClaimCommandHandler struct {
	log    logger.Logger
	tracer tracing.Tracer
	jwt    ose_jwt.Manager
}

// Handle implements cqrs.CommandHandle.
func (u parseClaimCommandHandler) Handle(ctx context.Context, command user.TokenCommand) (*ose_jwt.Claims, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.parse_claim.command.handler", trace.WithAttributes(
		attribute.String("operation", "parse_claim"),
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
			zap.String("operation", "parse_claim"),
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
			zap.String("operation", "parse_claim"),
			zap.Error(err),
		)

		return nil, err
	}

	return claims, nil
}

func newParseClaimCommandHandler(log logger.Logger, tracer tracing.Tracer,
	jwt ose_jwt.Manager) cqrs.CommandHandle[user.TokenCommand, *ose_jwt.Claims] {
	return &parseClaimCommandHandler{
		log:    log,
		tracer: tracer,
		jwt:    jwt,
	}
}

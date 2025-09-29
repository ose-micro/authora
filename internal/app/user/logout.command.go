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

type logoutCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	cache  user.Cache
	tracer tracing.Tracer
	bs     business.Domain
	jwt    ose_jwt.Manager
}

// Handle implements cqrs.CommandHandle.
func (h logoutCommandHandler) Handle(ctx context.Context, command user.TokenCommand) (bool, error) {
	ctx, span := h.tracer.Start(ctx, "app.user.logout.command.handler", trace.WithAttributes(
		attribute.String("operation", "logout"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command dto
	if err := command.Validate(); err != nil {
		err := ose_error.New(ose_error.ErrBadRequest, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "logout"),
			zap.Error(err),
		)

		return false, err
	}

	if err := h.cache.Delete(ctx, command.Token); err != nil {
		err = ose_error.New(ose_error.ErrUnauthorized, "failed to logout")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "logout"),
			zap.Error(err),
		)
		return false, err
	}

	h.log.Info("logout process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "logout"),
		zap.Any("dto", command),
	)

	return true, nil
}

func newLogoutCommandHandler(log logger.Logger, tracer tracing.Tracer, cache user.Cache) cqrs.CommandHandle[user.TokenCommand, bool] {
	return &logoutCommandHandler{
		log:    log,
		tracer: tracer,
		cache:  cache,
	}
}

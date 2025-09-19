package permission

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/permission"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type readOneQueryHandler struct {
	repo   permission.Repo
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readOneQueryHandler) Handle(ctx context.Context, query permission.ReadQuery) (*permission.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "app.permission.repository.query.handler", trace.WithAttributes(
		attribute.String("operation", "repository"),
		attribute.String("dto", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch permissions from store
	record, err := r.repo.ReadOne(ctx, query.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository permissions",
			zap.String("trace_id", traceId),
			zap.String("operation", "repository"),
			zap.Error(err),
		)

		return nil, err
	}

	r.log.Info("repository process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "repository"),
		zap.Any("dto", fmt.Sprintf("%v", query)),
	)
	return record, nil
}

func newReadOneQueryHandler(repo permission.Repo, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[permission.ReadQuery, *permission.Domain] {
	return &readOneQueryHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
	}
}

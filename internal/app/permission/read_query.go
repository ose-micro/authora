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

type readQueryHandler struct {
	repo   permission.Repo
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readQueryHandler) Handle(ctx context.Context, query permission.ReadQuery) (map[string]any, error) {
	ctx, span := r.tracer.Start(ctx, "app.permission.read.query.handler", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("dto", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch permissions from store
	records, err := r.repo.Read(ctx, query.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read permissions",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)

		return nil, err
	}

	r.log.Info("read process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "read"),
		zap.Any("dto", fmt.Sprintf("%v", query)),
	)
	return records, nil
}

func newReadQueryHandler(repo permission.Repo, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[permission.ReadQuery, map[string]any] {
	return &readQueryHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
	}
}

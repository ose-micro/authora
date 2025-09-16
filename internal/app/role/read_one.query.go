package role

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type readOneQueryHandler struct {
	repo   role.Repo
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readOneQueryHandler) Handle(ctx context.Context, query role.ReadQuery) (*role.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "app.role.read_one.query.handler", trace.WithAttributes(
		attribute.String("operation", "read_one"),
		attribute.String("dto", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch roles from store
	record, err := r.repo.ReadOne(ctx, query.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read roles",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)

		return nil, err
	}

	r.log.Info("read_one process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "read_one"),
		zap.Any("dto", fmt.Sprintf("%v", query)),
	)

	return record, nil
}

func newReadOneQueryHandler(repo role.Repo, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[role.ReadQuery, *role.Domain] {
	return &readOneQueryHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
	}
}

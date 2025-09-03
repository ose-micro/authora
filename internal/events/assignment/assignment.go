package assignment

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type defaultEvent struct {
	app    assignment.App
	log    logger.Logger
	tracer tracing.Tracer
}

func (d defaultEvent) AssignUserRole(ctx context.Context, event user.DefaultEvent) (*assignment.Domain, error) {
	ctx, span := d.tracer.Start(ctx, "events.assignment.created.handler", trace.WithAttributes(
		attribute.String("operation", "created"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	create, err := d.app.Create(ctx, assignment.CreateCommand{
		User:   event.ID,
		Tenant: event.Tenant,
		Role:   event.Role,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to process events",
			zap.String("trace_id", traceId),
			zap.String("operation", "created"),
			zap.Error(err),
		)
		return nil, err
	}

	return create, nil
}

func NewEvent(app assignment.App, log logger.Logger, tracer tracing.Tracer) assignment.Event {
	return defaultEvent{
		app:    app,
		log:    log,
		tracer: tracer,
	}
}

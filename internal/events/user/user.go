package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type defaultEvent struct {
	app    user.App
	log    logger.Logger
	tracer tracing.Tracer
}

func (d defaultEvent) ChangeStatus(ctx context.Context, event user.DefaultEvent) (bool, error) {
	ctx, span := d.tracer.Start(ctx, "events.user.change_state.handler", trace.WithAttributes(
		attribute.String("operation", "change_state"),
		attribute.String("dto", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	create, err := d.app.ChangeStatus(ctx, user.StatusCommand{
		Id:    event.ID,
		State: event.Status.State,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to process events",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_state"),
			zap.Error(err),
		)
		return false, err
	}

	return create, nil
}

func (d defaultEvent) OnBoard(ctx context.Context, event user.DefaultEvent) (*user.Domain, error) {
	ctx, span := d.tracer.Start(ctx, "events.user.onboard.handler", trace.WithAttributes(
		attribute.String("operation", "onboard"),
		attribute.String("dto", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	create, err := d.app.Create(ctx, user.CreateCommand{
		GivenNames: event.GivenNames,
		FamilyName: event.FamilyName,
		Email:      event.Email,
		Role:       event.Role,
		Password:   event.Password,
		Metadata:   event.Metadata,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to process events",
			zap.String("trace_id", traceId),
			zap.String("operation", "onboard"),
			zap.Error(err),
		)
		return nil, err
	}

	return create, nil
}

func NewEvent(app user.App, log logger.Logger, tracer tracing.Tracer) user.Event {
	return defaultEvent{
		app:    app,
		log:    log,
		tracer: tracer,
	}
}

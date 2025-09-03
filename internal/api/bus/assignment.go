package bus

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/events"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs/bus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func newAssignmentConsumer(bus bus.Bus, event events.Events, tracer tracing.Tracer, log logger.Logger) error {
	return bus.Subscribe(user.CreatedEvent, "assignment", func(ctx context.Context, data any) error {
		ctx, span := tracer.Start(ctx, "event.assignment.created.handler", trace.WithAttributes(
			attribute.String("operation", "created"),
			attribute.String("payload", fmt.Sprintf("%v", data)),
		))
		defer span.End()

		traceId := trace.SpanContextFromContext(ctx).TraceID().String()

		msg, err := toUserEvent(data)
		if err != nil {
			return err
		}
		_, err = event.Assignment.AssignUserRole(ctx, msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to broadcast event",
				zap.Any("msg", msg),
				zap.String("trace_id", traceId),
				zap.String("operation", "created"),
				zap.Error(err),
			)

			return err
		}

		return nil
	})
}

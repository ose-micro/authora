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

func newUserConsumer(bus bus.Bus, event events.Events, tracer tracing.Tracer, log logger.Logger) error {
	err := bus.Subscribe(user.CreatedEvent, "user_created_consumer", func(ctx context.Context, data any) error {
		ctx, span := tracer.Start(ctx, "event.user.created.handler", trace.WithAttributes(
			attribute.String("operation", "user_created"),
			attribute.String("payload", fmt.Sprintf("%v", data)),
		))
		defer span.End()

		traceId := trace.SpanContextFromContext(ctx).TraceID().String()

		msg, err := toUserEvent(data)

		if err != nil {
			return err
		}
		if _, err := event.User.ChangeStatus(ctx, msg); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to broadcast event",
				zap.Any("msg", msg),
				zap.String("trace_id", traceId),
				zap.String("operation", "user_created"),
				zap.Error(err),
			)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func newUserChangeStateEvent(bus bus.Bus, event events.Events, tracer tracing.Tracer, log logger.Logger) error {
	if err := bus.Subscribe(user.ChangeStateEvent, "user_change_state_consumer", func(ctx context.Context, data any) error {
		ctx, span := tracer.Start(ctx, "event.user.change_state.handler", trace.WithAttributes(
			attribute.String("operation", "change_state"),
			attribute.String("payload", fmt.Sprintf("%v", data)),
		))
		defer span.End()

		traceId := trace.SpanContextFromContext(ctx).TraceID().String()

		msg, err := toUserEvent(data)
		if err != nil {
			return err
		}
		if _, err = event.User.ChangeStatus(ctx, msg); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to broadcast event",
				zap.Any("msg", msg),
				zap.String("trace_id", traceId),
				zap.String("operation", "change_state"),
				zap.Error(err),
			)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

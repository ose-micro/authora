package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/events"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func newAssignmentConsumer(bus domain.Bus, app app.Apps, event events.Events, tracer tracing.Tracer, log logger.Logger) error {
	_ = bus.Subscribe(user.CreatedEvent, "assignment_user_create_consumer", func(ctx context.Context, data any) error {
		ctx, span := tracer.Start(ctx, "event.assignment.created.handler", trace.WithAttributes(
			attribute.String("operation", "created"),
			attribute.String("dto", fmt.Sprintf("%v", data)),
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

	_ = bus.Subscribe(assignment.AssignmentOnboardEvent, "assignment_onboard", func(ctx context.Context, data any) error {
		ctx, span := tracer.Start(ctx, "event.assignment.onboard.handler", trace.WithAttributes(
			attribute.String("operation", "onboard"),
			attribute.String("dto", fmt.Sprintf("%v", data)),
		))
		defer span.End()

		traceId := trace.SpanContextFromContext(ctx).TraceID().String()

		msg, err := toAssignmentEvent(data)
		if err != nil {
			return err
		}
		_, err = app.Assignment.Create(ctx, assignment.CreateCommand{
			User:   msg.User,
			Tenant: msg.Tenant,
			Role:   msg.Role,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to broadcast event",
				zap.Any("msg", msg),
				zap.String("trace_id", traceId),
				zap.String("operation", "onboard"),
				zap.Error(err),
			)

			return err
		}

		return nil
	})

	return nil
}

func toAssignmentEvent(data interface{}) (assignment.DefaultEvent, error) {
	mapData, ok := data.(map[string]interface{})
	if !ok {
		return assignment.DefaultEvent{}, fmt.Errorf("invalid message format")
	}

	// Marshal it to JSON
	raw, err := json.Marshal(mapData)
	if err != nil {
		return assignment.DefaultEvent{}, fmt.Errorf("failed to marshal map: %w", err)
	}

	var event assignment.DefaultEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		return event, fmt.Errorf("failed to unmarshal into DefaultEvent: %w", err)
	}

	return event, nil
}

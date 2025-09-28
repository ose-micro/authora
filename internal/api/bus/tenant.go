package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ose-micro/authora/internal/business/tenant"
	"github.com/ose-micro/authora/internal/events"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func newTenantConsumer(bus domain.Bus, event events.Events, tracer tracing.Tracer, log logger.Logger) error {
	return bus.Subscribe(tenant.NewEvent, "tenant-consumer", "new-tenant-consumer", func(ctx context.Context, data any) error {
		ctx, span := tracer.Start(ctx, "bus.tenant.onboard.handler", trace.WithAttributes(
			attribute.String("operation", "onboard"),
			attribute.String("dto", fmt.Sprintf("%v", data)),
		))
		defer span.End()

		traceId := trace.SpanContextFromContext(ctx).TraceID().String()

		msg, err := toTenantEvent(data)
		if err != nil {
			return err
		}

		_, err = event.Tenant.OnBoard(ctx, msg)
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
}

func toTenantEvent(data interface{}) (tenant.OnboardEvent, error) {
	mapData, ok := data.(map[string]interface{})
	if !ok {
		return tenant.OnboardEvent{}, fmt.Errorf("invalid message format")
	}

	// Marshal it to JSON
	raw, err := json.Marshal(mapData)
	if err != nil {
		return tenant.OnboardEvent{}, fmt.Errorf("failed to marshal map: %w", err)
	}

	var event tenant.OnboardEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		return event, fmt.Errorf("failed to unmarshal into Campaign: %w", err)
	}

	return event, nil
}

package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/tenant"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/events"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InvokeConsumers(lc fx.Lifecycle, app app.Apps, event *events.Events, bus domain.Bus, trancer tracing.Tracer, log logger.Logger) error {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				eventList := []string{
					tenant.NewEvent,
					assignment.AssignmentOnboardEvent,
					user.CreatedEvent,
					user.ChangeStateEvent,
				}
				err := bus.EnsureStream("EVENT", eventList...)
				if err != nil {
					log.Fatal("nats stream failed", zap.Error(err))
				}

				err = newTenantConsumer(bus, *event, trancer, log)
				if err != nil {
					log.Fatal("tenant consumer failed", zap.Error(err))
				}

				err = newUserConsumer(bus, *event, trancer, log)
				if err != nil {
					log.Fatal("user consumer failed", zap.Error(err))
				}
				
				err = newAssignmentConsumer(bus, app, *event, trancer, log)
				if err != nil {
					log.Fatal("assignment consumer failed", zap.Error(err))
				}
			}()

			return nil
		},
	})

	return nil
}

func toUserEvent(data interface{}) (user.DefaultEvent, error) {
	mapData, ok := data.(map[string]interface{})
	if !ok {
		return user.DefaultEvent{}, fmt.Errorf("invalid message format")
	}

	// Marshal it to JSON
	raw, err := json.Marshal(mapData)
	if err != nil {
		return user.DefaultEvent{}, fmt.Errorf("failed to marshal map: %w", err)
	}

	var event user.DefaultEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		return event, fmt.Errorf("failed to unmarshal into User: %w", err)
	}

	return event, nil
}

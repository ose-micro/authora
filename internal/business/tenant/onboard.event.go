package tenant

import (
	"time"

	"github.com/ose-micro/core/domain"
)

type OnboardEvent struct {
	Name      string                 `json:"name"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}

func (e OnboardEvent) EventName() string {
	return OnboardedEvent
}

func (e OnboardEvent) OccurredAt() time.Time {
	return e.CreatedAt
}

var _ domain.Event = OnboardEvent{}

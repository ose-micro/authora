package user

import (
	"time"

	"github.com/ose-micro/core/domain"
)

type DefaultEvent struct {
	ID         string                 `json:"_id"`
	GivenNames string                 `json:"given_names"`
	FamilyName string                 `json:"family_name"`
	Email      string                 `json:"email"`
	Password   string                 `json:"password"`
	Metadata   map[string]interface{} `json:"metadata"`
	Role       string                 `json:"role"`
	Tenant     string                 `json:"tenant"`
	Status     Status                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
}

func (e DefaultEvent) EventName() string {
	return "default_event"
}

func (e DefaultEvent) OccurredAt() time.Time {
	return e.CreatedAt
}

var _ domain.Event = DefaultEvent{}

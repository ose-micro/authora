package tenant

import (
	"context"
)

type Event interface {
	OnBoard(ctx context.Context, event OnboardEvent) (*Domain, error)
}

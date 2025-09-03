package user

import (
	"context"
)

type Event interface {
	OnBoard(ctx context.Context, event DefaultEvent) (*Domain, error)
}

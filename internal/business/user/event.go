package user

import (
	"context"
)

type Event interface {
	OnBoard(ctx context.Context, event DefaultEvent) (*Domain, error)
	ChangeStatus(ctx context.Context, event DefaultEvent) (bool, error)
}

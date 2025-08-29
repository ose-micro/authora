package tenant

import (
	"context"
)

type App interface {
	Create(ctx context.Context, command CreateCommand) (*Domain, error)
	Update(ctx context.Context, command UpdateCommand) (*Domain, error)
	Delete(ctx context.Context, command UpdateCommand) (*Domain, error)
	Read(ctx context.Context, command ReadQuery) (map[string]any, error)
}

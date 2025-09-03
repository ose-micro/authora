package assignment

import "context"

type App interface {
	Create(ctx context.Context, params CreateCommand) (*Domain, error)
	Update(ctx context.Context, params UpdateCommand) (*Domain, error)
	Delete(ctx context.Context, params UpdateCommand) (*Domain, error)
	Read(ctx context.Context, command ReadQuery) (map[string]any, error)
}

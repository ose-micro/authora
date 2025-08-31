package user

import "context"

type App interface {
	Create(ctx context.Context, params CreateCommand) (*Domain, error)
	Update(ctx context.Context, params UpdateCommand) (*Domain, error)
	ChangePassword(ctx context.Context, params ChangePasswordCommand) (*Domain, error)
	Login(ctx context.Context, params LoginCommand) (*Auth, error)
	Delete(ctx context.Context, params UpdateCommand) (*Domain, error)
	Read(ctx context.Context, command ReadQuery) (map[string]any, error)
}

package user

import (
	"context"

	ose_jwt "github.com/ose-micro/jwt"
)

type App interface {
	Create(ctx context.Context, command CreateCommand) (*Domain, error)
	Update(ctx context.Context, command UpdateCommand) (*Domain, error)
	ChangePassword(ctx context.Context, command ChangePasswordCommand) (*Domain, error)
	Login(ctx context.Context, command LoginCommand) (*Auth, error)
	Delete(ctx context.Context, command UpdateCommand) (*Domain, error)
	HasRole(ctx context.Context, command HasRoleCommand) (bool, error)
	HasPermission(ctx context.Context, command HasPermissionCommand) (bool, error)
	ParseClaims(ctx context.Context, command TokenCommand) (*ose_jwt.Claims, error)
	Read(ctx context.Context, command ReadQuery) (map[string]any, error)
}

package assignment

import (
	"context"

	"github.com/ose-micro/authora/internal/business/user"
)

type Event interface {
	AssignUserRole(ctx context.Context, event user.DefaultEvent) (*Domain, error)
}

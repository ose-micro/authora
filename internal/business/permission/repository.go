package permission

import (
	"context"

	"github.com/ose-micro/core/dto"
	ose_error "github.com/ose-micro/error"
)

type Repo interface {
	Create(ctx context.Context, payload Domain) *ose_error.Error
	Read(ctx context.Context, request dto.Request) (map[string]any, *ose_error.Error)
	ReadOne(ctx context.Context, request dto.Request) (*Domain, *ose_error.Error)
	Update(ctx context.Context, payload Domain) *ose_error.Error
	Delete(ctx context.Context, payload Domain) *ose_error.Error
}

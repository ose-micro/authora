package tenant

import (
	"time"

	"github.com/ose-micro/authora/internal/domain/tenant"
	"github.com/ose-micro/core/domain"
)

type Tenant struct {
	Id        string                 `bson:"_id"`
	Name      string                 `bson:"name"`
	Metadata  map[string]interface{} `bson:"metadata"`
	Version   int32                  `bson:"version"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
	DeletedAt *time.Time             `bson:"deleted_at"`
	Events    []domain.Event         `bson:"events"`
}

func newCollection(params tenant.Domain) Tenant {
	return Tenant{
		Id:        params.ID(),
		Name:      params.Name(),
		Metadata:  params.Metadata(),
		Version:   params.Version(),
		CreatedAt: params.CreatedAt(),
		UpdatedAt: params.UpdatedAt(),
		DeletedAt: params.DeletedAt(),
		Events:    params.Events(),
	}
}

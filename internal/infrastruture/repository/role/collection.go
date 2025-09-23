package role

import (
	"time"

	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/core/domain"
)

type Role struct {
	Id          string         `bson:"_id"`
	Name        string         `bson:"name"`
	Tenant      string         `bson:"tenant"`
	Description string         `bson:"description"`
	Permissions []string       `bson:"permissions"`
	Version     int32          `bson:"version"`
	CreatedAt   time.Time      `bson:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at"`
	DeletedAt   *time.Time     `bson:"deleted_at"`
	Events      []domain.Event `bson:"events"`
}

func newCollection(params role.Domain) Role {
	return Role{
		Id:          params.ID(),
		Name:        params.Name(),
		Tenant:      params.Tenant(),
		Permissions: params.Permissions(),
		Description: params.Description(),
		Version:     params.Version(),
		CreatedAt:   params.CreatedAt(),
		UpdatedAt:   params.UpdatedAt(),
		DeletedAt:   params.DeletedAt(),
		Events:      params.Events(),
	}
}

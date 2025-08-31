package permission

import (
	"time"

	"github.com/ose-micro/authora/internal/domain/permission"
	"github.com/ose-micro/core/domain"
)

type Permission struct {
	Id        string         `bson:"_id"`
	Resource  string         `bson:"resource"`
	Action    string         `bson:"action"`
	Version   int32          `bson:"version"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at"`
	DeletedAt *time.Time     `bson:"deleted_at"`
	Events    []domain.Event `bson:"events"`
}

func newCollection(params permission.Domain) Permission {
	return Permission{
		Id:        params.ID(),
		Resource:  params.Resource(),
		Action:    params.Action(),
		Version:   params.Version(),
		CreatedAt: params.CreatedAt(),
		UpdatedAt: params.UpdatedAt(),
		DeletedAt: params.DeletedAt(),
		Events:    params.Events(),
	}
}

package assignment

import (
	"time"

	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/core/domain"
)

type User struct {
	Id        string         `bson:"_id"`
	User      string         `bson:"user"`
	Tenant    string         `bson:"tenant"`
	Role      string         `bson:"role"`
	Version   int32          `bson:"version"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at"`
	DeletedAt *time.Time     `bson:"deleted_at"`
	Events    []domain.Event `bson:"events"`
}

func newCollection(params assignment.Domain) User {
	return User{
		Id:        params.ID(),
		User:      params.User(),
		Tenant:    params.Tenant(),
		Role:      params.Role(),
		Version:   params.Version(),
		CreatedAt: params.CreatedAt(),
		UpdatedAt: params.UpdatedAt(),
		DeletedAt: params.DeletedAt(),
		Events:    params.Events(),
	}
}

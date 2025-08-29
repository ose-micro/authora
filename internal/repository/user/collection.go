package user

import (
	"time"

	"github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/core/domain"
)

type User struct {
	Id         string                 `bson:"_id"`
	GivenNames string                 `bson:"given_names"`
	FamilyName string                 `bson:"family_name"`
	Email      string                 `bson:"email"`
	Password   string                 `bson:"password"`
	Metadata   map[string]interface{} `bson:"metadata"`
	Version    int32                  `bson:"version"`
	CreatedAt  time.Time              `bson:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at"`
	DeletedAt  *time.Time             `bson:"deleted_at"`
	Events     []domain.Event         `bson:"events"`
}

func newCollection(params user.Domain) User {
	return User{
		Id:         params.ID(),
		GivenNames: params.GivenNames(),
		FamilyName: params.FamilyName(),
		Email:      params.Email(),
		Password:   params.Password(),
		Metadata:   params.Metadata(),
		Version:    params.Version(),
		CreatedAt:  params.CreatedAt(),
		UpdatedAt:  params.UpdatedAt(),
		DeletedAt:  params.DeletedAt(),
		Events:     params.Events(),
	}
}

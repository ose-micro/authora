package user

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type Domain struct {
	*domain.Aggregate
	givenNames string
	familyName string
	email      string
	metadata   map[string]interface{}
}

type Params struct {
	Aggregate  *domain.Aggregate
	GivenNames string
	FamilyName string
	Email      string
	metadata   map[string]interface{}
}

type Public struct {
	Id         string                 `json:"_id"`
	GivenNames string                 `json:"given_names"`
	FamilyName string                 `json:"family_name"`
	Email      string                 `json:"email"`
	Metadata   map[string]interface{} `json:"metadata"`
	Version    int32                  `json:"version"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	DeletedAt  *time.Time             `json:"deleted_at"`
	Events     []domain.Event         `json:"events"`
}

func (d *Domain) GivenNames() string {
	return d.givenNames
}

func (d *Domain) FamilyName() string {
	return d.familyName
}

func (d *Domain) Email() string {
	return d.email
}

func (d *Domain) Metadata() map[string]interface{} {
	return d.metadata
}

func (d *Domain) Update(params Params) {
	if params.metadata != nil {
		d.metadata = params.metadata
		d.Touch()
	}

	if params.GivenNames != "" {
		d.givenNames = params.GivenNames
		d.Touch()
	}

	if params.FamilyName != "" {
		d.familyName = params.FamilyName
		d.Touch()
	}
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:         d.ID(),
		GivenNames: d.givenNames,
		FamilyName: d.familyName,
		Email:      d.email,
		Metadata:   d.metadata,
		Version:    d.Version(),
		CreatedAt:  d.CreatedAt(),
		UpdatedAt:  d.UpdatedAt(),
		DeletedAt:  d.DeletedAt(),
		Events:     d.Events(),
	}
}

func (p Public) Params() *Params {
	id := rid.Existing(p.Id)
	version := p.Version
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt
	deletedAt := p.DeletedAt
	events := p.Events

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Params{
		Aggregate:  aggregate,
		GivenNames: p.GivenNames,
		FamilyName: p.FamilyName,
		Email:      p.Email,
		metadata:   p.Metadata,
	}
}

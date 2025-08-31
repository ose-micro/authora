package role

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type Domain struct {
	*domain.Aggregate
	name        string
	tenant      string
	description string
	permissions []string
}

type Params struct {
	Aggregate   *domain.Aggregate
	Name        string   `json:"name"`
	Tenant      string   `json:"tenant"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type Public struct {
	Id          string         `json:"_id"`
	Name        string         `json:"name"`
	Tenant      string         `json:"tenant"`
	Permissions []string       `json:"permissions"`
	Description string         `json:"description"`
	Version     int32          `json:"version"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at"`
	Events      []domain.Event `json:"events"`
}

func (d *Domain) Name() string {
	return d.name
}

func (d *Domain) Tenant() string {
	return d.tenant
}

func (d *Domain) Description() string {
	return d.description
}

func (d *Domain) Permissions() []string {
	return d.permissions
}

func (d *Domain) Equals(other Domain) bool {
	return d.ID() == other.ID() && d.Version() == other.Version()
}

func (d *Domain) Update(params Params) {
	if params.Name != d.Name() {
		d.name = params.Name
		d.Touch()
	}

	if params.Tenant != d.Tenant() {
		d.tenant = params.Tenant
		d.Touch()
	}

	if params.Description != d.Description() {
		d.description = params.Description
		d.Touch()
	}

	if params.Permissions != nil {
		d.permissions = params.Permissions
		d.Touch()
	}
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:          d.ID(),
		Name:        d.name,
		Tenant:      d.tenant,
		Permissions: d.permissions,
		Description: d.description,
		Version:     d.Version(),
		CreatedAt:   d.CreatedAt(),
		UpdatedAt:   d.UpdatedAt(),
		DeletedAt:   d.DeletedAt(),
		Events:      d.Events(),
	}
}

func (p *Public) Params() *Params {
	id := rid.Existing(p.Id)
	version := p.Version
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt
	deletedAt := p.DeletedAt
	events := p.Events

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Params{
		Aggregate:   aggregate,
		Name:        p.Name,
		Tenant:      p.Tenant,
		Permissions: p.Permissions,
		Description: p.Description,
	}
}

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
	permissions []string
}

type Params struct {
	Aggregate   *domain.Aggregate
	Name        string   `json:"name"`
	Tenant      string   `json:"tenant"`
	Permissions []string `json:"permissions"`
}

type Public struct {
	Id          string         `json:"_id"`
	Name        string         `json:"name"`
	Tenant      string         `json:"tenant"`
	Permissions []string       `json:"permissions"`
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

func (d *Domain) Permissions() []string {
	return d.permissions
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:          d.ID(),
		Name:        d.name,
		Tenant:      d.tenant,
		Permissions: d.permissions,
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
	}
}

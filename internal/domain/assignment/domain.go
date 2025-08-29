package assignment

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type Domain struct {
	*domain.Aggregate
	user   string
	tenant string
	roles  []string
}

type Params struct {
	Aggregate *domain.Aggregate
	User      string
	Tenant    string
	Roles     []string
}

type Public struct {
	Id        string         `json:"_id"`
	User      string         `json:"user"`
	Tenant    string         `json:"tenant"`
	Roles     []string       `json:"roles"`
	Version   int32          `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	Events    []domain.Event `json:"events"`
}

func (d *Domain) User() string {
	return d.user
}

func (d *Domain) Tenant() string {
	return d.tenant
}

func (d *Domain) Roles() []string {
	return d.roles
}

func (d *Domain) UpdateRole(roles []string) {
	d.roles = roles
	d.Touch()
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:        d.ID(),
		User:      d.user,
		Tenant:    d.tenant,
		Roles:     d.roles,
		Version:   d.Version(),
		CreatedAt: d.CreatedAt(),
		UpdatedAt: d.UpdatedAt(),
		DeletedAt: d.DeletedAt(),
		Events:    d.Events(),
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
		Aggregate: aggregate,
		User:      p.User,
		Tenant:    p.Tenant,
		Roles:     p.Roles,
	}
}

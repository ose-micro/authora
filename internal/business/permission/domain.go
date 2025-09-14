package permission

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type Domain struct {
	*domain.Aggregate
	resource string
	action   string
}

type Params struct {
	Aggregate *domain.Aggregate
	Resource  string
	Action    string
}

type Public struct {
	Id        string         `json:"_id"`
	Resource  string         `json:"resource"`
	Action    string         `json:"action"`
	Version   int32          `json:"version"`
	Count     int32          `json:"count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	Events    []domain.Event `json:"events"`
}

func (d *Domain) Resource() string {
	return d.resource
}

func (d *Domain) Action() string {
	return d.action
}

func (d *Domain) Update(params Params) {
	if params.Resource != "" {
		d.resource = params.Resource
		d.Touch()
	}

	if params.Action != "" {
		d.action = params.Action
		d.Touch()
	}
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:        d.ID(),
		Resource:  d.resource,
		Action:    d.action,
		CreatedAt: d.CreatedAt(),
		UpdatedAt: d.UpdatedAt(),
		DeletedAt: d.DeletedAt(),
		Events:    d.Events(),
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
		Aggregate: aggregate,
		Resource:  p.Resource,
		Action:    p.Action,
	}
}

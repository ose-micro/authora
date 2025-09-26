package tenant

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

const NewEvent string = "authora.tenant_new"

type Domain struct {
	*domain.Aggregate
	name     string
	metadata map[string]interface{}
}

type Params struct {
	Aggregate *domain.Aggregate
	Name      string
	Metadata  map[string]interface{}
}

type Public struct {
	Id        string                 `json:"_id"`
	Name      string                 `json:"name"`
	Metadata  map[string]interface{} `json:"metadata"`
	Count     int32                  `json:"count"`
	Version   int32                  `json:"version"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at"`
	Events    []domain.Event         `json:"events"`
}

func (d *Domain) Name() string {
	return d.name
}

func (d *Domain) Metadata() map[string]interface{} {
	return d.metadata
}

func (d *Domain) Update(params Params) {
	if params.Metadata != nil {
		d.metadata = params.Metadata
		d.Touch()
	}

	if params.Name != "" {
		d.name = params.Name
		d.Touch()
	}
}

func (d *Domain) Public() *Public {
	return &Public{
		Id:        d.ID(),
		Name:      d.name,
		Metadata:  d.metadata,
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
		Name:      p.Name,
		Metadata:  p.Metadata,
	}
}

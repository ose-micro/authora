package role

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type initializer struct {
}

func (i initializer) New(param Params) (*Domain, error) {
	id := rid.New("rol", true)

	aggregate := domain.NewAggregate(*id)

	return &Domain{
		Aggregate:   aggregate,
		name:        param.Name,
		tenant:      param.Tenant,
		permissions: param.Permissions,
	}, nil
}

func (i initializer) Existing(param Params) (*Domain, error) {
	id := rid.Existing(param.Aggregate.ID())
	version := param.Aggregate.Version()
	createdAt := param.Aggregate.CreatedAt()
	updatedAt := param.Aggregate.UpdatedAt()
	deletedAt := param.Aggregate.DeletedAt()
	events := param.Aggregate.Events()

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Domain{
		Aggregate:   aggregate,
		name:        param.Name,
		tenant:      param.Tenant,
		permissions: param.Permissions,
	}, nil
}

func NewDomain() domain.Domain[Domain, Params] {
	return &initializer{}
}

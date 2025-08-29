package user

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/utils"
	ose_error "github.com/ose-micro/error"
	"github.com/ose-micro/rid"
)

type initializer struct {
}

func (i initializer) New(param Params) (*Domain, error) {
	id := rid.New("usr", true)

	aggregate := domain.NewAggregate(*id)
	password, err := utils.HashPassword(param.Password)
	if err != nil {
		return nil, ose_error.New(ose_error.ErrInvalidInput, err.Error())
	}

	return &Domain{
		Aggregate:  aggregate,
		givenNames: param.GivenNames,
		familyName: param.FamilyName,
		email:      param.Email,
		metadata:   param.Metadata,
		password:   password,
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
		Aggregate:  aggregate,
		givenNames: param.GivenNames,
		familyName: param.FamilyName,
		email:      param.Email,
		metadata:   param.Metadata,
	}, nil
}

func NewDomain() domain.Domain[Domain, Params] {
	return &initializer{}
}

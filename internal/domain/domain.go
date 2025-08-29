package domain

import (
	"github.com/ose-micro/authora/internal/domain/assignment"
	"github.com/ose-micro/authora/internal/domain/role"
	"github.com/ose-micro/authora/internal/domain/tenant"
	"github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/core/domain"
)

type Domain struct {
	Tenant     domain.Domain[tenant.Domain, tenant.Params]
	Role       domain.Domain[role.Domain, role.Params]
	User       domain.Domain[user.Domain, user.Params]
	Assignment domain.Domain[assignment.Domain, assignment.Params]
}

func Inject() Domain {
	return Domain{
		Tenant:     tenant.NewDomain(),
		Role:       role.NewDomain(),
		User:       user.NewDomain(),
		Assignment: assignment.NewDomain(),
	}
}

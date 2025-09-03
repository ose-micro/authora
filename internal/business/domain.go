package business

import (
	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/permission"
	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/authora/internal/business/tenant"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/core/domain"
)

type Domain struct {
	Tenant     domain.Domain[tenant.Domain, tenant.Params]
	Role       domain.Domain[role.Domain, role.Params]
	User       domain.Domain[user.Domain, user.Params]
	Permission domain.Domain[permission.Domain, permission.Params]
	Assignment domain.Domain[assignment.Domain, assignment.Params]
}

func Inject() Domain {
	return Domain{
		Tenant:     tenant.NewDomain(),
		Role:       role.NewDomain(),
		User:       user.NewDomain(),
		Assignment: assignment.NewDomain(),
		Permission: permission.NewDomain(),
	}
}

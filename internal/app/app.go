package app

import (
	"github.com/ose-micro/authora/internal/app/assignment"
	"github.com/ose-micro/authora/internal/app/permission"
	"github.com/ose-micro/authora/internal/app/role"
	"github.com/ose-micro/authora/internal/app/tenant"
	"github.com/ose-micro/authora/internal/app/user"
	"github.com/ose-micro/authora/internal/business"
	assignmentDomain "github.com/ose-micro/authora/internal/business/assignment"
	permissionDomain "github.com/ose-micro/authora/internal/business/permission"
	roleDomain "github.com/ose-micro/authora/internal/business/role"
	tenantDomain "github.com/ose-micro/authora/internal/business/tenant"
	userDomain "github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_jwt "github.com/ose-micro/jwt"
)

type Apps struct {
	Tenant     tenantDomain.App
	Role       roleDomain.App
	User       userDomain.App
	Assignment assignmentDomain.App
	Permission permissionDomain.App
}

func Inject(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer, manager *ose_jwt.Manager, bus domain.Bus) Apps {

	return Apps{
		Tenant:     tenant.NewApp(bs, log, tracer, repo),
		Role:       role.NewApp(bs, log, tracer, repo),
		User:       user.NewApp(bs, log, tracer, repo, *manager, bus),
		Assignment: assignment.NewApp(bs, log, tracer, repo),
		Permission: permission.NewApp(bs, log, tracer, repo),
	}
}

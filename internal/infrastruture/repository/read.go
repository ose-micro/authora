package repository

import (
	"github.com/ose-micro/authora/internal/business"
	assignmentDomain "github.com/ose-micro/authora/internal/business/assignment"
	permissionDomain "github.com/ose-micro/authora/internal/business/permission"
	roleDomain "github.com/ose-micro/authora/internal/business/role"
	tenantDomain "github.com/ose-micro/authora/internal/business/tenant"
	userDomain "github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/infrastruture/repository/assignment"
	"github.com/ose-micro/authora/internal/infrastruture/repository/permission"
	"github.com/ose-micro/authora/internal/infrastruture/repository/role"
	"github.com/ose-micro/authora/internal/infrastruture/repository/tenant"
	"github.com/ose-micro/authora/internal/infrastruture/repository/user"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	mongodb "github.com/ose-micro/mongo"
)

type Repository struct {
	Tenant     tenantDomain.Repo
	Role       roleDomain.Repo
	User       userDomain.Repo
	Assignment assignmentDomain.Repo
	Permission permissionDomain.Repo
}

func Inject(db *mongodb.Client, bs business.Domain, log logger.Logger, tracer tracing.Tracer) Repository {
	return Repository{
		Tenant:     tenant.NewRepository(db, log, tracer, bs),
		Role:       role.NewRepository(db, log, tracer, bs),
		User:       user.NewRepository(db, log, tracer, bs),
		Assignment: assignment.NewRepository(db, log, tracer, bs),
		Permission: permission.NewRepository(db, log, tracer, bs),
	}
}

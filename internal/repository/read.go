package repository

import (
	"github.com/ose-micro/authora/internal/domain"
	assignmentDomain "github.com/ose-micro/authora/internal/domain/assignment"
	permissionDomain "github.com/ose-micro/authora/internal/domain/permission"
	roleDomain "github.com/ose-micro/authora/internal/domain/role"
	tenantDomain "github.com/ose-micro/authora/internal/domain/tenant"
	userDomain "github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/authora/internal/repository/assignment"
	"github.com/ose-micro/authora/internal/repository/permission"
	"github.com/ose-micro/authora/internal/repository/role"
	"github.com/ose-micro/authora/internal/repository/tenant"
	"github.com/ose-micro/authora/internal/repository/user"
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

func Inject(db *mongodb.Client, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) Repository {
	return Repository{
		Tenant:     tenant.NewRepository(db, log, tracer, bs),
		Role:       role.NewRepository(db, log, tracer, bs),
		User:       user.NewRepository(db, log, tracer, bs),
		Assignment: assignment.NewRepository(db, log, tracer, bs),
		Permission: permission.NewRepository(db, log, tracer, bs),
	}
}

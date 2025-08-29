package app

import (
	"github.com/ose-micro/authora/internal/app/assignment"
	"github.com/ose-micro/authora/internal/app/role"
	"github.com/ose-micro/authora/internal/app/tenant"
	"github.com/ose-micro/authora/internal/app/user"
	"github.com/ose-micro/authora/internal/domain"
	assignmentDomain "github.com/ose-micro/authora/internal/domain/assignment"
	roleDomain "github.com/ose-micro/authora/internal/domain/role"
	tenantDomain "github.com/ose-micro/authora/internal/domain/tenant"
	userDomain "github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
)

type Apps struct {
	Tenant     tenantDomain.App
	Role       roleDomain.App
	User       userDomain.App
	Assignment assignmentDomain.App
}

func Inject(bs domain.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer) Apps {

	return Apps{
		Tenant:     tenant.NewApp(bs, log, tracer, repo),
		Role:       role.NewApp(bs, log, tracer, repo),
		User:       user.NewApp(bs, log, tracer, repo),
		Assignment: assignment.NewApp(bs, log, tracer, repo),
	}
}

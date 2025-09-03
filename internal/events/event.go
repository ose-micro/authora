package events

import (
	"github.com/ose-micro/authora/internal/app"
	assignmentBusiness "github.com/ose-micro/authora/internal/business/assignment"
	tenantBusiness "github.com/ose-micro/authora/internal/business/tenant"
	userBusiness "github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/events/assignment"
	"github.com/ose-micro/authora/internal/events/tenant"
	"github.com/ose-micro/authora/internal/events/user"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
)

type Events struct {
	Tenant     tenantBusiness.Event
	User       userBusiness.Event
	Assignment assignmentBusiness.Event
}

func NewEvents(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *Events {
	return &Events{
		Tenant:     tenant.NewEvent(apps.Tenant, log, tracer),
		User:       user.NewEvent(apps.User, log, tracer),
		Assignment: assignment.NewEvent(apps.Assignment, log, tracer),
	}
}

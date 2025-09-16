package assignment

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/domain"
)

type DefaultEvent struct {
	User   string
	Tenant string
	Role   string
}

func (c DefaultEvent) CommandName() string {
	return "assignment.create.command"
}

func (c DefaultEvent) Validate() error {
	fields := make([]string, 0)

	if c.User == "" {
		fields = append(fields, "user is required")
	}

	if c.Tenant == "" {
		fields = append(fields, "tenant is required")
	}

	if c.Role == "" {
		fields = append(fields, "role is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ domain.Command = DefaultEvent{}

package assignment

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type CreateCommand struct {
	User   string
	Tenant string
	Roles  []string
}

func (c CreateCommand) CommandName() string {
	return "assignment.create.command"
}

func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.User == "" {
		fields = append(fields, "user is required")
	}

	if c.Tenant == "" {
		fields = append(fields, "tenant is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}

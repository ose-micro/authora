package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type HasRoleCommand struct {
	Token  string
	Role   string
	Tenant string
}

func (c HasRoleCommand) CommandName() string {
	return "user.has.role.command"
}

func (c HasRoleCommand) Validate() error {
	fields := make([]string, 0)

	if c.Role == "" {
		fields = append(fields, "role is required")
	}

	if c.Tenant == "" {
		fields = append(fields, "tenant is required")
	}

	if c.Token == "" {
		fields = append(fields, "token is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = HasRoleCommand{}

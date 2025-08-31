package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/common"
	"github.com/ose-micro/cqrs"
)

type HasPermissionCommand struct {
	Token      string
	Tenant     string
	Permission *common.Permission
}

func (c HasPermissionCommand) CommandName() string {
	return "user.has.role.command"
}

func (c HasPermissionCommand) Validate() error {
	fields := make([]string, 0)

	if c.Tenant == "" {
		fields = append(fields, "tenant is required")
	}

	if c.Permission == nil {
		fields = append(fields, "permissions is required")
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

var _ cqrs.Command = HasPermissionCommand{}

package role

import (
	"errors"
	"strings"

	"github.com/ose-micro/common"
	"github.com/ose-micro/cqrs"
)

type CreateCommand struct {
	Name        string
	Tenant      string
	Description string
	Permissions []common.Permission
}

func (c CreateCommand) CommandName() string {
	return "role.create.command"
}

func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Name == "" {
		fields = append(fields, "name is required")
	}

	if c.Description == "" {
		fields = append(fields, "description is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}

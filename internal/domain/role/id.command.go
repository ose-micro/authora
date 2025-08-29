package role

import (
	"github.com/ose-micro/cqrs"
)

type IdCommand struct {
	Id string
}

func (c IdCommand) CommandName() string {
	return "role.id.command"
}

func (c IdCommand) Validate() error {
	fields := make([]string, 0)

	if c.Id == "" {
		fields = append(fields, "id is required")
	}

	return nil
}

var _ cqrs.Command = IdCommand{}

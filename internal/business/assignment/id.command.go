package assignment

import (
	"github.com/ose-micro/core/domain"
)

type IdCommand struct {
	Id string
}

func (c IdCommand) CommandName() string {
	return "user_id.id.command"
}

func (c IdCommand) Validate() error {
	fields := make([]string, 0)

	if c.Id == "" {
		fields = append(fields, "id is required")
	}

	return nil
}

var _ domain.Command = IdCommand{}

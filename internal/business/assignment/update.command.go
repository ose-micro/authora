package assignment

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/domain"
)

type UpdateCommand struct {
	Id   string
	Role string
}

func (u UpdateCommand) CommandName() string {
	return "assignment.update.command"
}

func (u UpdateCommand) Validate() error {
	fields := make([]string, 0)

	if u.Id == "" {
		fields = append(fields, "id is required")
	}

	if u.Role == "" {
		fields = append(fields, "role is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ domain.Command = UpdateCommand{}

package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/domain"
)

type StatusCommand struct {
	Id    string
	State State
}

func (c StatusCommand) CommandName() string {
	return "user.status.command"
}

func (c StatusCommand) Validate() error {
	fields := make([]string, 0)

	if c.Id == "" {
		fields = append(fields, "id is required")
	}

	if c.State == -1 {
		fields = append(fields, "status is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ domain.Command = StatusCommand{}

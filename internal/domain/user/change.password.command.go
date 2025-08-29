package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type ChangePasswordCommand struct {
	Id          string
	Password    string
	OldPassword string
}

func (c ChangePasswordCommand) CommandName() string {
	return "user.create.command"
}

func (c ChangePasswordCommand) Validate() error {
	fields := make([]string, 0)

	if c.OldPassword == "" {
		fields = append(fields, "old password is required")
	}

	if c.Password == "" {
		fields = append(fields, "password is required")
	} else if len(c.Password) < 8 {
		fields = append(fields, "password must be at least 8 characters")
	}

	if c.Id == "" {
		fields = append(fields, "id is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = ChangePasswordCommand{}

package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type ResetPasswordCommand struct {
	Id          string
	NewPassword string
}

func (c ResetPasswordCommand) CommandName() string {
	return "user.reset_password.command"
}

func (c ResetPasswordCommand) Validate() error {
	fields := make([]string, 0)

	if c.NewPassword == "" {
		fields = append(fields, "new password is required")
	} else if len(c.NewPassword) < 8 {
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

var _ cqrs.Command = ResetPasswordCommand{}

package user

import (
	"github.com/ose-micro/cqrs"
)

type LoginCommand struct {
	Email    string
	Password string
}

func (c LoginCommand) CommandName() string {
	return "login.id.command"
}

func (c LoginCommand) Validate() error {
	fields := make([]string, 0)

	if c.Email == "" {
		fields = append(fields, "email is required")
	}

	if c.Password == "" {
		fields = append(fields, "password is required")
	}

	return nil
}

var _ cqrs.Command = LoginCommand{}

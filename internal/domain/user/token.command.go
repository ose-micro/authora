package user

import (
	"github.com/ose-micro/cqrs"
)

type TokenCommand struct {
	Token string
}

func (c TokenCommand) CommandName() string {
	return "user.token.command"
}

func (c TokenCommand) Validate() error {
	fields := make([]string, 0)

	if c.Token == "" {
		fields = append(fields, "token is required")
	}

	return nil
}

var _ cqrs.Command = TokenCommand{}

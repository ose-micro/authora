package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/utils"
	"github.com/ose-micro/cqrs"
)

type CreateCommand struct {
	Name     string
	Email    string
	Metadata map[string]interface{}
}

func (c CreateCommand) CommandName() string {
	return "user.create.command"
}

func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Name == "" {
		fields = append(fields, "name is required")
	}

	if c.Email == "" {
		fields = append(fields, "email is required")
	} else if !utils.IsValidEmail(c.Email) {
		fields = append(fields, "email is invalid")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}

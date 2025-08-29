package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type UpdateCommand struct {
	Id         string
	GivenNames string
	FamilyName string
	Metadata   map[string]interface{}
}

func (u UpdateCommand) CommandName() string {
	return "user.update.command"
}

func (u UpdateCommand) Validate() error {
	fields := make([]string, 0)

	if u.Id == "" {
		fields = append(fields, "id is required")
	}

	if u.GivenNames == "" {
		fields = append(fields, "given names is required")
	}

	if u.FamilyName == "" {
		fields = append(fields, "family name is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = UpdateCommand{}

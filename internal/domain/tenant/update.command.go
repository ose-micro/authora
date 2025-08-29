package tenant

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type UpdateCommand struct {
	Id       string
	Name     string
	Metadata map[string]interface{}
}

func (u UpdateCommand) CommandName() string {
	return "update.command.tenant"
}

func (u UpdateCommand) Validate() error {
	fields := make([]string, 0)

	if u.Id == "" {
		fields = append(fields, "id is required")
	}

	if u.Name == "" {
		fields = append(fields, "name is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = UpdateCommand{}

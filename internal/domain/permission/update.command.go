package permission

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type UpdateCommand struct {
	Id       string
	Resource string
	Action   string
}

func (u UpdateCommand) CommandName() string {
	return "permission.update.command"
}

func (u UpdateCommand) Validate() error {
	fields := make([]string, 0)

	if u.Id == "" {
		fields = append(fields, "id is required")
	}

	if u.Resource == "" {
		fields = append(fields, "resource is required")
	}

	if u.Action == "" {
		fields = append(fields, "action is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = UpdateCommand{}

package permission

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/domain"
)

type CreateCommand struct {
	Resource string
	Action   string
}

func (c CreateCommand) CommandName() string {
	return "permission.create.command"
}

func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Resource == "" {
		fields = append(fields, "resource is required")
	}

	if c.Action == "" {
		fields = append(fields, "action is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ domain.Command = CreateCommand{}

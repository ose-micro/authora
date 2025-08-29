package tenant

import (
	"errors"
	"strings"

	"github.com/ose-micro/cqrs"
)

type CreateCommand struct {
	Name     string
	Metadata map[string]interface{}
}

func (c CreateCommand) CommandName() string {
	return "create.command.tenant"
}

func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Name == "" {
		fields = append(fields, "name is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}

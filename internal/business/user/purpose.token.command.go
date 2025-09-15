package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/domain"
)

type PurposeTokenCommand struct {
	Id      string
	Purpose string
	Safe    bool
}

func (c PurposeTokenCommand) CommandName() string {
	return "user.create.command"
}

func (c PurposeTokenCommand) Validate() error {
	fields := make([]string, 0)

	if c.Id == "" {
		fields = append(fields, "id is required")
	}

	if c.Purpose == "" {
		fields = append(fields, "purpose is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ domain.Command = PurposeTokenCommand{}

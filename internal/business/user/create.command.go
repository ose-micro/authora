package user

import (
	"errors"
	"strings"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/utils"
)

type CreateCommand struct {
	GivenNames string
	FamilyName string
	Email      string
	Password   string
	Role       string
	Tenant     string
	Metadata   map[string]interface{}
}

func (c CreateCommand) CommandName() string {
	return "user.create.command"
}

func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.GivenNames == "" {
		fields = append(fields, "given names is required")
	}

	if c.FamilyName == "" {
		fields = append(fields, "family name is required")
	}

	if c.Password == "" {
		fields = append(fields, "password is required")
	} else if len(c.Password) < 8 {
		fields = append(fields, "password must be at least 8 characters")
	}

	if c.Email == "" {
		fields = append(fields, "email is required")
	} else if !utils.IsValidEmail(c.Email) {
		fields = append(fields, "email is invalid")
	}

	if c.Role == "" {
		fields = append(fields, "role is required")
	}

	if c.Tenant == "" {
		fields = append(fields, "tenant is required")
	}

	if len(fields) > 0 {
		msg := strings.Join(fields, " ")
		return errors.New(msg)
	}

	return nil
}

var _ domain.Command = CreateCommand{}

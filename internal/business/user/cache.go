package user

import (
	"context"
	"errors"
	"strings"
	"time"

	ose_jwt "github.com/ose-micro/jwt"
	"github.com/ose-micro/rid"
)

type Token struct {
	key     string
	user    string
	tenant  string
	purpose string
	token   string
}

type TokenParam struct {
	Key     string `json:"key"`
	User    string `json:"user"`
	Purpose string `json:"purpose"`
	Tenant  string `json:"tenant"`
	Token   string `json:"token"`
}

type TokenClaim struct {
	Key     string         `json:"key"`
	User    string         `json:"user"`
	Purpose string         `json:"purpose"`
	Tenant  string         `json:"tenant"`
	Token   string         `json:"token"`
	Claims  ose_jwt.Claims `json:"claims"`
}

func (t *Token) Token() string {
	return t.token
}

func (t *Token) User() string {
	return t.user
}

func (t *Token) Tenant() string {
	return t.tenant
}

func (t *Token) Purpose() string {
	return t.purpose
}

func (t *Token) Key() string {
	return t.key
}

func (t *Token) Param() TokenParam {
	return TokenParam{
		Key:     t.key,
		User:    t.user,
		Purpose: t.purpose,
		Tenant:  t.tenant,
		Token:   t.token,
	}
}

func NewToken(param TokenParam) (*Token, error) {
	if err := param.Validate(false); err != nil {
		return nil, err
	}

	id := rid.New(param.Purpose, true)
	return &Token{
		key:     id.String(),
		user:    param.User,
		tenant:  param.Tenant,
		purpose: param.Purpose,
		token:   param.Token,
	}, nil
}

func ExistingToken(param TokenParam) (*Token, error) {
	if err := param.Validate(true); err != nil {
		return nil, err
	}

	return &Token{
		key:     param.Key,
		user:    param.User,
		tenant:  param.Tenant,
		purpose: param.Purpose,
		token:   param.Token,
	}, nil
}

type Cache interface {
	Save(ctx context.Context, payload *Token, ttl time.Duration) error
	Get(ctx context.Context, key string) (*Token, error)
}

func (p TokenParam) Validate(isExisting bool) error {
	msg := make([]string, 0)

	if p.Token == "" {
		msg = append(msg, "token is required")
	}

	if p.Purpose == "" {
		msg = append(msg, "purpose is required")
	}

	if p.Tenant == "" {
		msg = append(msg, "tenant is required")
	}

	if p.User == "" {
		msg = append(msg, "user is required")
	}

	if isExisting {
		if p.Key == "" {
			msg = append(msg, "key is required")
		}
	}

	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "; "))
	}

	return nil
}

package user

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/dto"
)

type ReadOneQuery struct {
	Request dto.Request
}

// QueryName implements cqrs.Query.
func (q ReadOneQuery) QueryName() string {
	return "user.read_one.query"
}

var _ domain.Query = ReadQuery{}

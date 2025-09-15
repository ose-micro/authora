package assignment

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/dto"
)

type ReadQuery struct {
	Request dto.Request
}

// QueryName implements cqrs.Query.
func (c ReadQuery) QueryName() string {
	return "user.repository.query"
}

var _ domain.Query = ReadQuery{}

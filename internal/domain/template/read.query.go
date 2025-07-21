package template

import (
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/cqrs"
)

type ReadQuery struct {
	Request dto.Request
}

// QueryName implements cqrs.Query.
func (c ReadQuery) QueryName() string {
	return READ_QUERY
}

var _ cqrs.Query = ReadQuery{}

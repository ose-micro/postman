package template

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/rid"
)

type roleDomain struct {
	timestamp timestamp.Timestamp
}

// Existing Load implements IDomain.
func (r roleDomain) Existing(param Params) (*Domain, error) {
	id := rid.Existing(param.Aggregate.ID())
	version := param.Aggregate.Version()
	createdAt := param.Aggregate.CreatedAt()
	updatedAt := param.Aggregate.UpdatedAt()
	deletedAt := param.Aggregate.DeletedAt()
	events := param.Aggregate.Events()

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Domain{
		Aggregate:    aggregate,
		subject:      param.Subject,
		content:      param.Content,
		placeholders: param.Placeholders,
	}, nil
}

// New implements IDomain.
func (r roleDomain) New(param Params) (*Domain, error) {
	id := rid.New("tmp", true)

	aggregate := domain.NewAggregate(*id)

	return &Domain{
		Aggregate:    aggregate,
		subject:      param.Subject,
		content:      param.Content,
		placeholders: param.Placeholders,
	}, nil
}

func NewTemplateDomain(timestamp timestamp.Timestamp) domain.Domain[Domain, Params] {
	return &roleDomain{
		timestamp: timestamp,
	}
}

package email

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/rid"
)

type emailDomain struct {
	timestamp timestamp.Timestamp
}

// Existing Load implements IDomain.
func (e emailDomain) Existing(param Params) (*Domain, error) {
	id := rid.Existing(param.Aggregate.ID())
	version := param.Aggregate.Version()
	createdAt := param.Aggregate.CreatedAt()
	updatedAt := param.Aggregate.UpdatedAt()
	deletedAt := param.Aggregate.DeletedAt()
	events := param.Aggregate.Events()

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Domain{
		Aggregate: aggregate,
		recipient: param.Recipient,
		sender:    param.Sender,
		from:      param.From,
		subject:   param.Subject,
		template:  param.Template,
		data:      param.Data,
		message:   param.Message,
		state:     param.State,
	}, nil
}

// New implements IDomain.
func (e emailDomain) New(param Params) (*Domain, error) {
	id := rid.New("eml", true)

	aggregate := domain.NewAggregate(*id)

	return &Domain{
		Aggregate: aggregate,
		recipient: param.Recipient,
		sender:    param.Sender,
		from:      param.From,
		subject:   param.Subject,
		template:  param.Template,
		data:      param.Data,
		message:   param.Message,
		state:     param.State,
	}, nil
}

func NewEmailDomain(timestamp timestamp.Timestamp) domain.Domain[Domain, Params] {
	return &emailDomain{
		timestamp: timestamp,
	}
}

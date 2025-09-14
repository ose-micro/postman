package email

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type State string

const (
	StateComplete State = "Complement"
	StateFailed   State = "Failed"
)

type Domain struct {
	*domain.Aggregate
	recipient string
	sender    string
	from      string
	subject   string
	template  string
	data      map[string]interface{}
	message   string
	state     State
}

type Public struct {
	Id        string                 `json:"_id"`
	Recipient string                 `json:"recipient"`
	Sender    string                 `json:"sender"`
	Subject   string                 `json:"subject"`
	Count     int32                  `json:"count"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
	Version   int32                  `json:"version"`
	State     State                  `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at"`
	Events    []domain.Event         `json:"events"`
}

type Params struct {
	*domain.Aggregate
	Recipient string
	Sender    string
	Subject   string
	Data      map[string]interface{}
	Template  string
	From      string
	Message   string
	State     State
}

func (p Public) Params() *Params {
	id := rid.Existing(p.Id)
	version := p.Version
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt
	deletedAt := p.DeletedAt
	events := p.Events

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Params{
		Aggregate: aggregate,
		Recipient: p.Recipient,
		Sender:    p.Sender,
		Subject:   p.Subject,
		Data:      p.Data,
		Template:  p.Template,
		From:      p.From,
		Message:   p.Message,
		State:     p.State,
	}
}

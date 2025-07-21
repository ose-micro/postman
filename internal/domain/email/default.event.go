package email

import (
	"time"

	"github.com/ose-micro/cqrs"
)

type DomainEvent struct {
	Id        string                 `json:"id"`
	Recipient string                 `json:"recipient"`
	Sender    string                 `json:"sender"`
	Subject   string                 `json:"subject"`
	Count     int32                  `json:"count"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
	State     State                  `json:"status"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// EventName implements cqrs.Event.
func (c DomainEvent) EventName() string {
	return DEFAULT_EVENT
}

var _ cqrs.Event = DomainEvent{}

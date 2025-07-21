package template

import (
	"time"

	"github.com/ose-micro/cqrs"
)

type DomainEvent struct {
	Id           string    `json:"id"`
	Content      string    `json:"content"`
	Subject      string    `json:"subject"`
	Count        int32       `json:"count"`
	Placeholders []string  `json:"placeholders"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// EventName implements cqrs.Event.
func (u DomainEvent) EventName() string {
	return UPDATED_COMMAND
}

var _ cqrs.Event = DomainEvent{}

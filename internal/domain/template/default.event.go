package template

import (
	"time"

	"github.com/ose-micro/cqrs"
)

type DefaultEvent struct {
	Id           string    `json:"id"`
	Content      string    `json:"content"`
	Subject      string    `json:"subject"`
	Placeholders []string  `json:"placeholders"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// EventName implements cqrs.Event.
func (u DefaultEvent) EventName() string {
	return UPDATED_COMMAND
}

var _ cqrs.Event = DefaultEvent{}
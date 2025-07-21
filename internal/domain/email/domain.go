package email

import (
	"time"

	"github.com/ose-micro/core/timestamp"
)

type State string

const (
	StateComplete State = "Complement"
	StateFailed   State = "Failed"
)

type Domain struct {
	timestamp timestamp.Timestamp
	id        string
	recipient string
	sender    string
	from      string
	subject   string
	template  string
	data      map[string]interface{}
	message   string
	state     State
	createdAt time.Time
	updatedAt time.Time
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
	State     State                  `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type Params struct {
	Id        string
	Recipient string
	Sender    string
	Subject   string
	Data      map[string]interface{}
	Template  string
	From      string
	Message   string
	State     State
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
	sender    *string
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
	Id        string                 `json:"id"`
	Recipient string                 `json:"recipient"`
	Sender    *string                `json:"sender"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
	State     State                  `json:"status"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type Params struct {
	Id        string
	Recipient string
	Sender    *string
	Subject   string
	Data      map[string]interface{}
	Template  string
	From      string
	Message   string
	State     State
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Domain) GetID() string {
	return d.id
}

func (d *Domain) GetRecipient() string {
	return d.recipient
}

func (d *Domain) GetData() map[string]interface{} {
	return d.data
}

func (d *Domain) GetSender() *string {
	return d.sender
}

func (d *Domain) GetFrom() string {
	return d.from
}

func (d *Domain) GetSubject() string {
	return d.subject
}

func (d *Domain) GetTemplate() string {
	return d.template
}

func (d *Domain) GetMessage() string {
	return d.message
}

func (d *Domain) GetState() State {
	return d.state
}

func (d *Domain) GetCreatedAt() time.Time {
	return d.createdAt
}

func (d *Domain) GetUpdatedAt() time.Time {
	return d.updatedAt
}

func (d *Domain) SetState(state State) {
	d.state = state
}

func (d *Domain) MakePublic() Public {
	return Public{
		Id:        d.id,
		Recipient: d.recipient,
		Sender:    d.sender,
		Subject:   d.subject,
		Data:      d.data,
		Template:  d.template,
		From:      d.from,
		State:     d.state,
		Message:   d.message,
		CreatedAt: d.createdAt,
		UpdatedAt: d.updatedAt,
	}
}

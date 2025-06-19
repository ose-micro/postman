package email

import (
	"time"

	"github.com/moriba-cloud/ose-postman/internal/domain/email"
)

type Email struct {
	Id        string      `bson:"_id"`
	Recipient string      `bson:"recipient,omitempty"`
	Sender    *string      `bson:"sender,omitempty"`
	Subject   string      `bson:"subject,omitempty"`
	Data      any         `bson:"data,omitempty"`
	Template  string      `bson:"template,omitempty"`
	From      string      `bson:"from,omitempty"`
	Message   string      `bson:"message,omitempty"`
	State     email.State `bson:"state,omitempty"`
	CreatedAt time.Time   `bson:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at"`
}

func newCollection(params email.Domain) Email {
	return Email{
		Id:        params.GetID(),
		Recipient: params.GetRecipient(),
		Sender:    params.GetSender(),
		Subject:   params.GetSubject(),
		Data:      params.GetData(),
		Template:  params.GetTemplate(),
		From:      params.GetFrom(),
		Message:   params.GetMessage(),
		State:     params.GetState(),
		CreatedAt: params.GetCreatedAt(),
		UpdatedAt: params.GetUpdatedAt(),
	}
}

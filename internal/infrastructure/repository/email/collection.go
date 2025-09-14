package email

import (
	"time"

	"github.com/ose-micro/postman/internal/business/email"
)

type Email struct {
	Id        string                 `bson:"_id"`
	Recipient string                 `bson:"recipient,omitempty"`
	Sender    string                 `bson:"sender,omitempty"`
	Subject   string                 `bson:"subject,omitempty"`
	Data      map[string]interface{} `bson:"data,omitempty"`
	Template  string                 `bson:"template,omitempty"`
	From      string                 `bson:"from,omitempty"`
	Message   string                 `bson:"message,omitempty"`
	State     email.State            `bson:"state,omitempty"`
	Version   int32                  `bson:"version,omitempty"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
	DeletedAt *time.Time             `bson:"deleted_at"`
}

func newCollection(params email.Domain) Email {
	return Email{
		Id:        params.ID(),
		Recipient: params.Recipient(),
		Sender:    params.Sender(),
		Subject:   params.Subject(),
		Data:      params.Data(),
		Template:  params.Template(),
		From:      params.From(),
		Message:   params.Message(),
		State:     params.State(),
		Version:   params.Version(),
		CreatedAt: params.CreatedAt(),
		UpdatedAt: params.UpdatedAt(),
		DeletedAt: params.DeletedAt(),
	}
}

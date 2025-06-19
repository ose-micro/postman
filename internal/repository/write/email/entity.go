package email

import (
	"encoding/json"
	"log"
	"time"

	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/postgres"
)

type Email struct {
	Id        string `gorm:"primaryKey"`
	Recipient string `gorm:"not null;index"`
	Sender    *string `gorm:"not null;index"`
	Subject   string `gorm:"not null;index"`
	Data      string `gorm:"type:text"`
	Template  string
	From      string
	Message   string `gorm:"type:text"`
	State     email.State
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EntityName implements postgres.Entity.
func (u Email) EntityName() string {
	return "emails"
}

func newEntity(params email.Domain) Email {
	jsonData, err := json.Marshal(params.GetData())
	if err != nil {
		log.Fatalln(err)
	}

	data := string(jsonData)
	
	return Email{
		Id:        params.GetID(),
		Recipient: params.GetRecipient(),
		Sender:    params.GetSender(),
		Subject:   params.GetSubject(),
		Data:      data,
		Template:  params.GetTemplate(),
		From:      params.GetFrom(),
		Message:   params.GetMessage(),
		State:     params.GetState(),
		CreatedAt: params.GetCreatedAt(),
		UpdatedAt: params.GetUpdatedAt(),
	}
}

var _ postgres.Entity = Email{}

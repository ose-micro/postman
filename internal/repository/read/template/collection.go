package template

import (
	"time"

	"github.com/moriba-cloud/ose-postman/internal/domain/template"
)

type Collection struct {
	Id           string    `bson:"_id"`
	Subject      string    `bson:"subject"`
	Content      string    `bson:"content"`
	Placeholders []string  `bson:"placeholders"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

func newCollection(params template.Domain) Collection {
	return Collection{
		Id: params.GetID(),
		Subject: params.GetSubject(),
		Content: params.GetContent(),
		Placeholders: params.GetPlaceholders(),
		CreatedAt: params.GetCreatedAt(),
		UpdatedAt: params.GetUpdatedAt(),
	}
}

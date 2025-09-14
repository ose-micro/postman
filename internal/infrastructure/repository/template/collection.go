package template

import (
	"time"

	"github.com/ose-micro/postman/internal/business/template"
)

type Collection struct {
	Id           string     `bson:"_id"`
	Subject      string     `bson:"subject"`
	Content      string     `bson:"content"`
	Placeholders []string   `bson:"placeholders"`
	Version      int32      `bson:"version"`
	CreatedAt    time.Time  `bson:"created_at"`
	UpdatedAt    time.Time  `bson:"updated_at"`
	DeletedAt    *time.Time `bson:"deleted_at"`
}

func newCollection(params template.Domain) Collection {
	return Collection{
		Id:           params.ID(),
		Subject:      params.Subject(),
		Content:      params.Content(),
		Placeholders: params.Placeholders(),
		Version:      params.Version(),
		CreatedAt:    params.CreatedAt(),
		UpdatedAt:    params.UpdatedAt(),
		DeletedAt:    params.DeletedAt(),
	}
}

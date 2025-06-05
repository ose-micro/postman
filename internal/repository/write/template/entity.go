package template

import (
	"time"

	"github.com/lib/pq"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/postgres"
)

type Template struct {
	ID           string `gorm:"primaryKey;type:text"`
	Subject      string
	Content      string
	Placeholders pq.StringArray `gorm:"type:text[]"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// EntityName implements postgres.Entity.
func (t Template) EntityName() string {
	return "templates"
}

func newEntity(params template.Domain) Template {
	return Template{
		ID:           params.GetID(),
		Subject:      params.GetSubject(),
		Content:      params.GetContent(),
		Placeholders: params.GetPlaceholders(),
		CreatedAt:    params.GetCreatedAt(),
		UpdatedAt:    params.GetUpdatedAt(),
	}
}

var _ postgres.Entity = Template{}

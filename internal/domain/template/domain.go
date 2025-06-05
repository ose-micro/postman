package template

import (
	"time"

	"github.com/ose-micro/core/timestamp"
)

type Domain struct {
	timestamp    timestamp.Timestamp
	id           string
	subject      string
	content      string
	placeholders []string
	createdAt    time.Time
	updatedAt    time.Time
}

type Public struct {
	Id           string    `json:"id"`
	Content      string    `json:"content"`
	Subject      string    `json:"subject"`
	Placeholders []string  `json:"placeholders"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Params struct {
	Id           string
	Content      string
	Subject      string
	Placeholders []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (d *Domain) GetID() string {
	return d.id
}

func (d *Domain) GetContent() string {
	return d.content
}

func (d *Domain) GetSubject() string {
	return d.subject
}

func (d *Domain) GetPlaceholders() []string {
	return d.placeholders
}

func (d *Domain) GetCreatedAt() time.Time {
	return d.createdAt
}

func (d *Domain) GetUpdatedAt() time.Time {
	return d.updatedAt
}

func (d *Domain) Update(param Params) error {
	if param.Content != "" {
		d.content = param.Content
		d.placeholders = param.Placeholders
		d.updatedAt = d.timestamp.Now()
	}

	if param.Subject != "" {
		d.subject = param.Subject
		d.updatedAt = d.timestamp.Now()
	}

	return nil
}

func (d *Domain) MakePublic() Public {
	return Public{
		Id:           d.id,
		Content:      d.content,
		Subject:      d.subject,
		Placeholders: d.placeholders,
		
		CreatedAt:    d.createdAt,
		UpdatedAt:    d.updatedAt,
	}
}

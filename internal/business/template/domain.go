package template

import (
	"time"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/rid"
)

type Domain struct {
	*domain.Aggregate
	subject      string
	content      string
	placeholders []string
}

type Public struct {
	Id           string         `json:"_id"`
	Content      string         `json:"content"`
	Subject      string         `json:"subject"`
	Count        int32          `json:"count"`
	Placeholders []string       `json:"placeholders"`
	Version      int32          `json:"version"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at"`
	Events       []domain.Event `json:"events"`
}

type Params struct {
	*domain.Aggregate
	Content      string
	Subject      string
	Placeholders []string
}

func (p Public) Params() *Params {
	id := rid.Existing(p.Id)
	version := p.Version
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt
	deletedAt := p.DeletedAt
	events := p.Events

	aggregate := domain.ExistingAggregate(*id, version, createdAt, updatedAt, deletedAt, events)

	return &Params{
		Aggregate:    aggregate,
		Subject:      p.Subject,
		Content:      p.Content,
		Placeholders: p.Placeholders,
	}
}

func (d *Domain) Content() string {
	return d.content
}

func (d *Domain) Subject() string {
	return d.subject
}

func (d *Domain) Placeholders() []string {
	return d.placeholders
}

func (d *Domain) SetContent(content string) *Domain {
	if content != "" {
		d.content = content
		d.Touch()
	}

	return d
}

func (d *Domain) SetPlaceholders(placeholders []string) *Domain {
	d.placeholders = placeholders
	d.Touch()

	return d
}

func (d *Domain) SetSubject(subject string) *Domain {
	if subject != "" {
		d.content = subject
		d.Touch()
	}

	return d
}

func (d *Domain) Public() Public {
	return Public{
		Id:           d.ID(),
		Content:      d.content,
		Subject:      d.subject,
		Placeholders: d.placeholders,
		Version:      d.Version(),
		CreatedAt:    d.CreatedAt(),
		UpdatedAt:    d.UpdatedAt(),
		DeletedAt:    d.DeletedAt(),
	}
}

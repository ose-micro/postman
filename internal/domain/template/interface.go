package template

import (
	"context"

	"github.com/ose-micro/core/dto"
)

const (
	CREATED_COMMAND string = "template.created.command"
	UPDATED_COMMAND string = "template.updated.command"
	DELETED_COMMAND string = "template.deleted.command"

	QUEUE string = "template_queue"

	SEND_MAIL_EVENT string = "ose-postman.email.send_mail.event"

	READ_QUERY     string = "template.read.query"
	READ_ONE_QUERY string = "template.read_one.query"
)

type Write interface {
	Create(ctx context.Context, payload Domain) error
	Read(ctx context.Context, request dto.Query) (*Domain, error)
	Update(ctx context.Context, payload Domain) error
	Delete(ctx context.Context, payload Domain) error
}

type Read interface {
	Create(ctx context.Context, payload Domain) error
	Read(ctx context.Context, request dto.Request) (map[string]any, error)
	Update(ctx context.Context, payload Domain) error
	Delete(ctx context.Context, payload Domain) error
}

type App interface {
	Create(ctx context.Context, command CreateCommand) (Domain, error)
	Read(ctx context.Context, request dto.Request) (map[string]any, error)
	Update(ctx context.Context, command UpdateCommand) error
	Delete(ctx context.Context, command DeleteCommand) error
}

type Event interface {
	Created(ctx context.Context, event DomainEvent) error
	Updated(ctx context.Context, event DomainEvent) error
	Deleted(ctx context.Context, event DomainEvent) error
}

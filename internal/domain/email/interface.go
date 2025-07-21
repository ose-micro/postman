package email

import (
	"context"

	"github.com/ose-micro/core/dto"
)

const (
	CREATED_COMMAND string = "email.created.command"
	UPDATED_COMMAND string = "email.updated.command"
	DELETED_COMMAND string = "email.deleted.command"

	READ_QUERY     string = "email.read.query"
	READ_ONE_QUERY string = "email.read_one.query"

	QUEUE string = "email_queue"

	DEFAULT_EVENT string = ""

	SEND_MAIL_EVENT string = "ose-postman.email.send-mail.event"
)

type Write interface {
	Create(ctx context.Context, payload Domain) error
	Read(ctx context.Context, request dto.Query) (*Domain, error)
	Update(ctx context.Context, payload Domain) error
}

type Read interface {
	Create(ctx context.Context, payload Domain) error
	Read(ctx context.Context, request dto.Request) (map[string]any, error)
	Update(ctx context.Context, payload Domain) error
}

type App interface {
	Create(ctx context.Context, command CreateCommand) (*Domain, error)
	Read(ctx context.Context, request dto.Request) (map[string]any, error)
	Resend(ctx context.Context, command IdCommand) error
}

type Event interface {
	Created(ctx context.Context, event DomainEvent) error
	Updated(ctx context.Context, event DomainEvent) error
}

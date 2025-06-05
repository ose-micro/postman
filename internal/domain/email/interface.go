package email

import (
	"context"

	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/cqrs"
)

const (
	CREATED_COMMAND string = "email.created.command"
	UPDATED_COMMAND string = "email.updated.command"
	DELETED_COMMAND string = "email.deleted.command"
	READ_QUERY      string = "email.read.query"
	READ_ONE_QUERY  string = "email.read_one.query"

	SEND_MAIL_EVENT string = "ose-postman.email.send-mail.event"
)

type Repository interface {
	Create(ctx context.Context, payload Domain) error
	Read(ctx context.Context, request dto.Request) ([]Domain, error)
	ReadOne(ctx context.Context, filters ...dto.Filter) (*Domain, error)
	Update(ctx context.Context, payload Domain) error
	Delete(ctx context.Context, payload Domain) error
}

type App interface {
	Create(ctx context.Context, command cqrs.Command) (Domain, error)
	Read(ctx context.Context, request dto.Request) ([]Domain, error)
	Resend(ctx context.Context, command cqrs.Command) (Domain, error)
	ReadOne(ctx context.Context, filters ...dto.Filter) (Domain, error)
	Delete(ctx context.Context, command cqrs.Command) error
}

type Event interface {
	Created(ctx context.Context, event cqrs.Event) error
	Updated(ctx context.Context, event cqrs.Event) error
	Deleted(ctx context.Context, event cqrs.Event) error
	SendMail(ctx context.Context, event cqrs.Event) error
}

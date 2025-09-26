package email

import (
	"context"

	"github.com/ose-micro/core/dto"
)

const (
	SendMailEvent string = "postman.send_mail"
)

type Repo interface {
	Create(ctx context.Context, payload Domain) error
	Read(ctx context.Context, request dto.Request) (map[string]any, error)
	ReadOne(ctx context.Context, request dto.Request) (*Domain, error)
	Update(ctx context.Context, payload Domain) error
}

type App interface {
	Create(ctx context.Context, command CreateCommand) (*Domain, error)
	Read(ctx context.Context, request dto.Request) (map[string]any, error)
	Resend(ctx context.Context, command IdCommand) error
}

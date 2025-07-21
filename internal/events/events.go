package events

import (
	"github.com/moriba-cloud/ose-postman/internal/app"
	"github.com/moriba-cloud/ose-postman/internal/domain"
	emailDomain "github.com/moriba-cloud/ose-postman/internal/domain/email"
	templateDomain "github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/moriba-cloud/ose-postman/internal/events/email"
	"github.com/moriba-cloud/ose-postman/internal/events/template"
	"github.com/moriba-cloud/ose-postman/internal/repository/read"
	"github.com/moriba-cloud/ose-postman/internal/repository/write"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
)

type Events struct {
	Email    emailDomain.Event
	Template templateDomain.Event
}

func Inject(bs domain.Domain, repo read.Repository, log logger.Logger,
	tracer tracing.Tracer, app app.Apps, write write.Repository) Events {
	return Events{
		Email:    email.New(bs, repo.Email, log, tracer, app.Email),
		Template: template.New(bs, repo.Template, log, tracer, app.Template),
	}
}

package app

import (
	"github.com/moriba-cloud/ose-postman/internal/app/email"
	"github.com/moriba-cloud/ose-postman/internal/app/template"
	"github.com/moriba-cloud/ose-postman/internal/domain"
	emailDomain "github.com/moriba-cloud/ose-postman/internal/domain/email"
	domain_template "github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/moriba-cloud/ose-postman/internal/repository/read"
	"github.com/moriba-cloud/ose-postman/internal/repository/write"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs/bus"
	"github.com/ose-micro/mailer"
)

type Apps struct {
	Template domain_template.App
	Email    emailDomain.App
}

func InjectApps(bs domain.Domain, write write.Repository, read read.Repository, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus, mailer *mailer.Mailer) Apps {
	return Apps{
		Template: template.NewTemplateApp(bs, log, tracer, write.Template, read.Template, bus, mailer),
		Email:    email.NewEmailApp(bs, log, tracer, write, read.Email, bus, mailer),
	}
}

package app

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/mailer"
	"github.com/ose-micro/postman/internal/app/email"
	"github.com/ose-micro/postman/internal/app/template"
	"github.com/ose-micro/postman/internal/business"
	emailDomain "github.com/ose-micro/postman/internal/business/email"
	domain_template "github.com/ose-micro/postman/internal/business/template"
	"github.com/ose-micro/postman/internal/infrastructure/repository"
)

type Apps struct {
	Template domain_template.App
	Email    emailDomain.App
}

func InjectApps(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer, bus domain.Bus, mailer *mailer.Mailer) Apps {
	return Apps{
		Template: template.NewTemplateApp(bs, log, tracer, repo.Template, bus, mailer),
		Email:    email.NewEmailApp(bs, log, tracer, repo, bus, mailer),
	}
}

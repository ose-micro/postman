package app

import (
	"github.com/moriba-cloud/ose-postman/internal/app/email"
	"github.com/moriba-cloud/ose-postman/internal/app/template"
	"github.com/moriba-cloud/ose-postman/internal/domain"
	domain_email "github.com/moriba-cloud/ose-postman/internal/domain/email"
	domain_template "github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/moriba-cloud/ose-postman/internal/repository/read"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
)

type Events struct {
	Template domain_template.Event
	Email    domain_email.Event
}

func InjectEvents(bs domain.Domain, repo read.Repository, app Apps, log logger.Logger, tracer tracing.Tracer) Events {
	return Events{
		Template: template.NewTemplateEvent(bs, repo.Template, log, tracer),
		Email:    email.NewEmailEvent(bs, repo.Email, log, tracer, app.Email),
	}
}

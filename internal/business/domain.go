package business

import (
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/postman/internal/business/email"
	"github.com/ose-micro/postman/internal/business/template"
)

type Domain struct {
	Template domain.Domain[template.Domain, template.Params]
	Email    domain.Domain[email.Domain, email.Params]
}

func InjectDomain(timestamp timestamp.Timestamp) Domain {
	return Domain{
		Template: template.NewTemplateDomain(timestamp),
		Email:    email.NewEmailDomain(timestamp),
	}
}

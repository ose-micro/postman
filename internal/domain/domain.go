package domain

import (
	"github.com/moriba-cloud/ose-postman/internal/common"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/timestamp"
)

type Domain struct {
	Template common.IDomain[template.Domain, template.Params]
	Email    common.IDomain[email.Domain, email.Params]
}

func InjectDomain(timestamp timestamp.Timestamp) Domain {
	return Domain{
		Template: template.NewTemplateDomain(timestamp),
		Email:    email.NewEmailDomain(timestamp),
	}
}

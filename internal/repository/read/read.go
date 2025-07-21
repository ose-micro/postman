package read

import (
	"github.com/moriba-cloud/ose-postman/internal/domain"
	emailDomain "github.com/moriba-cloud/ose-postman/internal/domain/email"
	templateDomain "github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/moriba-cloud/ose-postman/internal/repository/read/email"
	"github.com/moriba-cloud/ose-postman/internal/repository/read/template"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	mongodb "github.com/ose-micro/mongo"
)

type Repository struct {
	Template templateDomain.Read
	Email    emailDomain.Read
}

func InjectRepository(db *mongodb.Client, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) Repository {
	return Repository{
		Template: template.NewRepository(db, log, tracer, bs),
		Email:    email.NewRepository(db, log, tracer, bs),
	}
}

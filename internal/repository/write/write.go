package write

import (
	"github.com/moriba-cloud/ose-postman/internal/domain"
	emailDomain "github.com/moriba-cloud/ose-postman/internal/domain/email"
	templateDomain "github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/moriba-cloud/ose-postman/internal/repository/write/email"
	"github.com/moriba-cloud/ose-postman/internal/repository/write/template"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/postgres"
)

type Repository struct {
	Template templateDomain.Repository
	Email    emailDomain.Repository
}

func InjectRepository(db *postgres.Postgres, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) Repository {
	return Repository{
		Template: template.NewTemplateRepository(db, bs, log, tracer),
		Email:    email.NewEmailRepository(db, bs, log, tracer),
	}
}

func Migrate(db *postgres.Postgres) error {
	return db.Conn().AutoMigrate(&template.Template{}, &email.Email{})
}

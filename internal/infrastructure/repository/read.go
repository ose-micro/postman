package repository

import (
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	mongodb "github.com/ose-micro/mongo"
	"github.com/ose-micro/postman/internal/business"
	emailDomain "github.com/ose-micro/postman/internal/business/email"
	templateDomain "github.com/ose-micro/postman/internal/business/template"
	"github.com/ose-micro/postman/internal/infrastructure/repository/email"
	"github.com/ose-micro/postman/internal/infrastructure/repository/template"
)

type Repository struct {
	Template templateDomain.Repo
	Email    emailDomain.Repo
}

func InjectRepository(db *mongodb.Client, bs business.Domain, log logger.Logger, tracer tracing.Tracer) Repository {
	return Repository{
		Template: template.NewRepository(db, log, tracer, bs),
		Email:    email.NewRepository(db, log, tracer, bs),
	}
}

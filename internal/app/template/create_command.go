package template

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	"github.com/ose-micro/mailer"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/template"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Handler
type createCommandHandler struct {
	repo   template.Repo
	log    logger.Logger
	mailer *mailer.Mailer
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (c *createCommandHandler) Handle(ctx context.Context, command template.CreateCommand) (*template.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.template.create.command.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process fail",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Any("details", err),
		)

		return nil, err
	}

	if err := c.mailer.ValidateData(command.Content, command.Placeholders); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process fail",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Any("details", err),
		)

		return nil, err
	}

	// create template
	domain, err := c.bs.Template.New(template.Params{
		Content:      command.Content,
		Subject:      command.Subject,
		Placeholders: command.Placeholders,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create template",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	// save template to write store
	err = c.repo.Create(ctx, *domain)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("fail while saving template",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	c.log.Info("create process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("payload", command),
	)
	return domain, nil
}

func newCreateCommandHandler(bs business.Domain, repo template.Repo, log logger.Logger,
	tracer tracing.Tracer, mailer *mailer.Mailer) cqrs.CommandHandle[template.CreateCommand, *template.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		mailer: mailer,
		bs:     bs,
	}
}

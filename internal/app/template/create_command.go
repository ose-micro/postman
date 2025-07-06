package template

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	"github.com/ose-micro/mailer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Handler
type createCommandHandler struct {
	repo   template.Repository
	log    logger.Logger
	mailer *mailer.Mailer
	bus    bus.Bus
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (c *createCommandHandler) Handle(ctx context.Context, command template.CreateCommand) (template.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.template.create.command.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process fail",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Any("details", err),
		)

		return template.Domain{}, err
	}

	if err := c.mailer.ValidateData(command.Content, command.Placeholders); err != nil {
		return template.Domain{}, err
	}

	// create domain
	domain, err := c.bs.Template.New(template.Params{
		Content:      command.Content,
		Subject:      command.Subject,
		Placeholders: command.Placeholders,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	// save template to write store
	err = c.repo.Create(ctx, *domain)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("fail while saving domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	// publish bus
	err = c.bus.Publish(command.CommandName(), template.DefaultEvent(domain.MakePublic()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("publish template created",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	c.log.Info("create process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.Any("payload", command),
	)
	return *domain, nil
}

func newCreateCommandHandler(bs domain.Domain, repo template.Repository, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus, mailer *mailer.Mailer) cqrs.CommandHandle[template.CreateCommand, template.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		mailer: mailer,
		bs:     bs,
	}
}

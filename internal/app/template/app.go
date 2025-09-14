package template

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	"github.com/ose-micro/mailer"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/template"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type templateApp struct {
	tracer tracing.Tracer
	log    logger.Logger
	create cqrs.CommandHandle[template.CreateCommand, *template.Domain]
	update cqrs.CommandHandle[template.UpdateCommand, bool]
	delete cqrs.CommandHandle[template.DeleteCommand, bool]
	read   cqrs.QueryHandle[template.ReadQuery, map[string]any]
}

// Create implements template.App.
func (t *templateApp) Create(ctx context.Context, command template.CreateCommand) (*template.Domain, error) {
	ctx, span := t.tracer.Start(ctx, "app.template.create.command", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := t.create.Handle(ctx, command)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	return record, nil
}

// Delete implements template.App.
func (t *templateApp) Delete(ctx context.Context, command template.DeleteCommand) error {
	ctx, span := t.tracer.Start(ctx, "app.template.delete.command", trace.WithAttributes(
		attribute.String("operation", "delete"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := t.delete.Handle(ctx, command); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "delete"),
			zap.Error(err),
		)

		return err
	}

	return nil
}

// Read implements template.App.
func (t *templateApp) Read(ctx context.Context, request dto.Request) (map[string]any, error) {
	ctx, span := t.tracer.Start(ctx, "app.email.read.query", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	records, err := t.read.Handle(ctx, template.ReadQuery{
		Request: request,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)
		return nil, err
	}

	return records, nil
}

// Update implements template.App.
func (t *templateApp) Update(ctx context.Context, command template.UpdateCommand) error {
	ctx, span := t.tracer.Start(ctx, "app.template.update.command", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := t.update.Handle(ctx, command); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func NewTemplateApp(bs business.Domain, log logger.Logger, tracer tracing.Tracer,
	repo template.Repo, bus bus.Bus, mailer *mailer.Mailer) template.App {
	return &templateApp{
		tracer: tracer,
		log:    log,
		create: newCreateCommandHandler(bs, repo, log, tracer, bus, mailer),
		read:   newReadQueryHandler(repo, log, tracer),
		update: newUpdateCommandHandler(bs, repo, log, tracer, bus),
		delete: newDeleteCommandHandler(bs, repo, log, tracer, bus),
	}
}

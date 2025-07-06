package template

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/dto"
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

type templateApp struct {
	tracer   tracing.Tracer
	log      logger.Logger
	create   cqrs.CommandHandle[template.CreateCommand, template.Domain]
	update   cqrs.CommandHandle[template.UpdateCommand, template.Domain]
	delete   cqrs.CommandHandle[template.DeleteCommand, template.Domain]
	read     cqrs.QueryHandle[ReadQuery, []template.Domain]
	read_one cqrs.QueryHandle[ReadOneQuery, template.Domain]
}

// Create implements template.App.
func (t *templateApp) Create(ctx context.Context, command cqrs.Command) (template.Domain, error) {
	ctx, span := t.tracer.Start(ctx, "app.template.create.command", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", command.(template.CreateCommand))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := t.create.Handle(ctx, command.(template.CreateCommand))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	return record, nil
}

// Delete implements template.App.
func (t *templateApp) Delete(ctx context.Context, command cqrs.Command) error {
	ctx, span := t.tracer.Start(ctx, "app.template.delete.command", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", command.(template.CreateCommand))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := t.delete.Handle(ctx, command.(template.DeleteCommand)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return err
	}

	return nil
}

// Read implements template.App.
func (t *templateApp) Read(ctx context.Context, request dto.Request) ([]template.Domain, error) {
	ctx, span := t.tracer.Start(ctx, "app.template.read.query", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	records, err := t.read.Handle(ctx, ReadQuery{
		request: request,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return nil, err
	}

	return records, nil
}

// ReadOne implements template.App.
func (t *templateApp) ReadOne(ctx context.Context, filters ...dto.Filter) (template.Domain, error) {
	ctx, span := t.tracer.Start(ctx, "app.template.read_one.query", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", filters)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := t.read_one.Handle(ctx, ReadOneQuery{
		filters: filters,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ_ONE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	return record, nil
}

// Update implements template.App.
func (t *templateApp) Update(ctx context.Context, command cqrs.Command) error {
	ctx, span := t.tracer.Start(ctx, "app.template.update.command", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", command.(template.UpdateCommand))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := t.update.Handle(ctx, command.(template.UpdateCommand)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func NewTemplateApp(bs domain.Domain, log logger.Logger, tracer tracing.Tracer, write template.Repository,
	read template.Repository, bus bus.Bus, mailer *mailer.Mailer) template.App {
	return &templateApp{
		tracer: tracer,
		log: log,
		create:   newCreateCommandHandler(bs, write, log, tracer, bus, mailer),
		read:     newReadQueryHandler(read, log, tracer),
		read_one: newReadOneQueryHandler(read, log, tracer),
		update:   newUpdateCommandHandler(bs, write, log, tracer, bus),
		delete:   newDeleteCommandHandler(bs, write, log, tracer, bus),
	}
}

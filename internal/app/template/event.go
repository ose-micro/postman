package template

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type templateEvent struct {
	log    logger.Logger
	tracer tracing.Tracer
	create cqrs.EventHandle[CreatedEvent]
	update cqrs.EventHandle[UpdatedEvent]
	delete cqrs.EventHandle[DeletedEvent]
}

// Created implements template.Event.
func (t *templateEvent) Created(ctx context.Context, event cqrs.Event) error {
	ctx, span := t.tracer.Start(ctx, "app.template.event.created", trace.WithAttributes(
		attribute.String("operation", "CREATED"),
		attribute.String("payload", fmt.Sprintf("%v", event.(CreatedEvent))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := t.create.Handle(ctx, event.(CreatedEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATED"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Deleted implements template.Event.
func (t *templateEvent) Deleted(ctx context.Context, event cqrs.Event) error {
	ctx, span := t.tracer.Start(ctx, "app.template.event.deleted", trace.WithAttributes(
		attribute.String("operation", "DELETED"),
		attribute.String("payload", fmt.Sprintf("%v", event.(DeletedEvent))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := t.delete.Handle(ctx, event.(DeletedEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETED"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Updated implements template.Event.
func (t *templateEvent) Updated(ctx context.Context, event cqrs.Event) error {
	ctx, span := t.tracer.Start(ctx, "app.template.event.updated", trace.WithAttributes(
		attribute.String("operation", "UPDATED"),
		attribute.String("payload", fmt.Sprintf("%v", event.(UpdatedEvent))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := t.update.Handle(ctx, event.(UpdatedEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATED"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func NewTemplateEvent(bs domain.Domain, repo template.Repository, log logger.Logger, tracer tracing.Tracer) template.Event {
	return &templateEvent{
		log: log,
		tracer: tracer,
		create: newCreatedEvent(bs, repo, log, tracer),
		update: newUpdatedEvent(bs, repo, log, tracer),
		delete: newDeletedEvent(bs, repo, log, tracer),
	}
}

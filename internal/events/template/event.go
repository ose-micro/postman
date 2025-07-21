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
	create cqrs.EventHandle[template.DomainEvent]
	update cqrs.EventHandle[template.DomainEvent]
}

// Updated implements template.Event.
func (r *templateEvent) Updated(ctx context.Context, event template.DomainEvent) error {
	ctx, span := r.tracer.Start(ctx, "app.template.updated.command", trace.WithAttributes(
		attribute.String("operation", "UPDATED"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	err := r.update.Handle(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATED"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Delete implements template.Event.
func (r *templateEvent) Deleted(ctx context.Context, event template.DomainEvent) error {
	ctx, span := r.tracer.Start(ctx, "app.template.delete.command", trace.WithAttributes(
		attribute.String("operation", "DELETED"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	err := r.update.Handle(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETED"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Credit implements template.Event.
func (r *templateEvent) Created(ctx context.Context, event template.DomainEvent) error {
	ctx, span := r.tracer.Start(ctx, "app.template.created.command", trace.WithAttributes(
		attribute.String("operation", "CREATED"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	err := r.create.Handle(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATED"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func New(bs domain.Domain, repo template.Read, log logger.Logger, tracer tracing.Tracer, app template.App) template.Event {
	return &templateEvent{
		log:      log,
		tracer:   tracer,
		create:   newCreatedEvent(bs, repo, log, tracer),
		update:   newUpdatedEvent(bs, repo, log, tracer),
	}
}

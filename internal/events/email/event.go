package email

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type emailEvent struct {
	log    logger.Logger
	tracer tracing.Tracer
	create cqrs.EventHandle[email.DomainEvent]
	update cqrs.EventHandle[email.DomainEvent]
}

// Updated implements email.Event.
func (r *emailEvent) Updated(ctx context.Context, event email.DomainEvent) error {
	ctx, span := r.tracer.Start(ctx, "app.email.updated.command", trace.WithAttributes(
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

// Credit implements email.Event.
func (r *emailEvent) Created(ctx context.Context, event email.DomainEvent) error {
	ctx, span := r.tracer.Start(ctx, "app.email.created.command", trace.WithAttributes(
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

func New(bs domain.Domain, repo email.Read, log logger.Logger, tracer tracing.Tracer, app email.App) email.Event {
	return &emailEvent{
		log:      log,
		tracer:   tracer,
		create:   newCreatedEvent(bs, repo, log, tracer),
		update:   newUpdatedEvent(bs, repo, log, tracer),
	}
}

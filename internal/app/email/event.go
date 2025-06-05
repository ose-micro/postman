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
	create cqrs.EventHandle[CreatedEvent]
	update cqrs.EventHandle[UpdatedEvent]
	delete cqrs.EventHandle[DeletedEvent]
	sendMail cqrs.EventHandle[SendMailEvent]
}

// Created implements email.Event.
func (e *emailEvent) Created(ctx context.Context, event cqrs.Event) error {
	ctx, span := e.tracer.Start(ctx, "app.email.event.created", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", event.(CreatedEvent))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := e.create.Handle(ctx, event.(CreatedEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Deleted implements email.Event.
func (e *emailEvent) Deleted(ctx context.Context, event cqrs.Event) error {
	ctx, span := e.tracer.Start(ctx, "app.email.event.deleted", trace.WithAttributes(
		attribute.String("operation", "DELETED"),
		attribute.String("payload", fmt.Sprintf("%v", event.(CreatedEvent))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := e.delete.Handle(ctx, event.(DeletedEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Updated implements email.Event.
func (e *emailEvent) Updated(ctx context.Context, event cqrs.Event) error {
	ctx, span := e.tracer.Start(ctx, "app.email.event.updated", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", event.(CreatedEvent))),
	))
	defer span.End()
	
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := e.update.Handle(ctx, event.(UpdatedEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	return nil
}


// SendMail implements email.Event.
func (e *emailEvent) SendMail(ctx context.Context, event cqrs.Event) error {
	ctx, span := e.tracer.Start(ctx, "app.email.event.send_mail", trace.WithAttributes(
		attribute.String("operation", "SEND_MAIL"),
		attribute.String("payload", fmt.Sprintf("%v", event.(SendMailEvent))),
	))
	defer span.End()
	
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if err := e.sendMail.Handle(ctx, event.(SendMailEvent)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process event",
			zap.String("trace_id", traceId),
			zap.String("operation", "SEND_MAIL"),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func NewEmailEvent(bs domain.Domain, repo email.Repository, log logger.Logger, tracer tracing.Tracer, app email.App) email.Event {
	return &emailEvent{
		log: log,
		tracer: tracer,
		create: newCreatedEvent(bs, repo, log, tracer),
		update: newUpdatedEvent(bs, repo, log, tracer),
		delete: newDeletedEvent(bs, repo, log, tracer),
		sendMail: newSendMailEvent(bs, app, log, tracer),
	}
}

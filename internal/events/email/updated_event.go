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

type updatedEventHandler struct {
	repo   email.Read
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

func (c *updatedEventHandler) ToDomain(event email.DomainEvent) (email.Domain, error) {
	record, err := c.bs.Email.Existing(email.Params{
		Id:        event.Id,
		Recipient: event.Recipient,
		Sender:    event.Sender,
		Subject:   event.Subject,
		Data:      event.Data,
		Template:  event.Template,
		From:      event.From,
		Message:   event.Message,
		State:     event.State,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	})

	if err != nil {
		return email.Domain{}, err
	}

	return *record, nil
}

// Handle implements cqrs.EventHandle.
func (c *updatedEventHandler) Handle(ctx context.Context, event email.DomainEvent) error {
	ctx, span := c.tracer.Start(ctx, "app.email.created.event.handler", trace.WithAttributes(
		attribute.String("operation", "UPDATED_EVENT"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// cast to domain
	payload, err := c.ToDomain(event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("fail to cast to domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATED_EVENT"),
			zap.Error(err),
		)

		return err
	}

	// save to db
	if err := c.repo.Update(ctx, payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to save to db",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATED_EVENT"),
			zap.Error(err),
		)

		return err
	}

	c.log.Info("created process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "UPDATED_EVENT"),
		zap.Any("payload", event),
	)
	return nil
}

func newUpdatedEvent(bs domain.Domain, repo email.Read, log logger.Logger,
	tracer tracing.Tracer) cqrs.EventHandle[email.DomainEvent] {
	return &updatedEventHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

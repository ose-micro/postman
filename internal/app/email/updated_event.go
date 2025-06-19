package email

import (
	"context"
	"fmt"
	"time"

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

type UpdatedEvent struct {
	Id        string                 `json:"id"`
	Recipient string                 `json:"recipient"`
	Sender    *string                `json:"sender"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
	State     email.State            `json:"status"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// EventName implements cqrs.Event.
func (c UpdatedEvent) EventName() string {
	return email.UPDATED_COMMAND
}

var _ cqrs.Event = UpdatedEvent{}

type updatedEvent struct {
	repo   email.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

func (c *updatedEvent) ToDomain(event UpdatedEvent) (*email.Domain, error) {
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
		return nil, err
	}

	return record, nil
}

// Handle implements cqrs.EventHandle.
func (c *updatedEvent) Handle(ctx context.Context, event UpdatedEvent) error {
	ctx, span := c.tracer.Start(ctx, "app.email.updated.event.handler", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// cast to domain
	payload, err := c.ToDomain(event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to cast to domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Any("details", err),
		)

		return err
	}

	// save to db
	if err := c.repo.Update(ctx, *payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to update to mongo",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)

		return err
	}

	c.log.Info("updated process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "UPDATE"),
		zap.Any("payload", event),
	)
	return nil
}

func newUpdatedEvent(bs domain.Domain, repo email.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.EventHandle[UpdatedEvent] {
	return &updatedEvent{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

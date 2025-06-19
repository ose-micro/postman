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

type CreatedEvent struct {
	Id        string                 `json:"id"`
	Recipient string                 `json:"recipient"`
	Sender    *string                `json:"sender"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
	State     email.State            `json:"state"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// EventName implements cqrs.Event.
func (c CreatedEvent) EventName() string {
	return "email.created.event"
}

var _ cqrs.Event = CreatedEvent{}

type createdEvent struct {
	repo   email.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

func (c *createdEvent) ToDomain(event CreatedEvent) (email.Domain, error) {
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
func (c *createdEvent) Handle(ctx context.Context, event CreatedEvent) error {
	ctx, span := c.tracer.Start(ctx, "app.email.created.event.handler", trace.WithAttributes(
		attribute.String("operation", "CREATED"),
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
			zap.String("stage", "app.handler"),
			zap.Any("details", err),
		)

		return err
	}

	// save to db
	if err := c.repo.Create(ctx, payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to save to db",
			zap.String("trace_id", traceId),
			zap.String("stage", "app.handler"),
			zap.Any("details", err),
		)

		return err
	}

	c.log.Info("created process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATED"),
		zap.Any("payload", event),
	)
	return nil
}

func newCreatedEvent(bs domain.Domain, repo email.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.EventHandle[CreatedEvent] {
	return &createdEvent{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

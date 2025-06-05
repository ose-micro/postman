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

type SendMailEvent struct {
	Recipient string                 `json:"recipient"`
	Sender    string                 `json:"sender"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
}

// EventName implements cqrs.Event.
func (s SendMailEvent) EventName() string {
	return email.SEND_MAIL_EVENT
}

var _ cqrs.Event = SendMailEvent{}

type sendMailEvent struct {
	app    email.App
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.EventHandle.
func (c *sendMailEvent) Handle(ctx context.Context, event SendMailEvent) error {
	ctx, span := c.tracer.Start(ctx, "app.email.send_mail.event.handler", trace.WithAttributes(
		attribute.String("operation", "SEND_MAIL"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// create mail
	if _, err := c.app.Create(ctx, CreateCommand{
		Recipient: event.Recipient,
		Sender:    event.Sender,
		Data:      event.Data,
		Template:  event.Template,
		From:      event.From,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create mail",
			zap.String("trace_id", traceId),
			zap.String("stage", "app.handler"),
			zap.Any("details", err),
		)

		return err
	}

	c.log.Info("send mail process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "SEND_MAIL"),
		zap.Any("payload", event),
	)
	return nil
}

func newSendMailEvent(bs domain.Domain, app email.App, log logger.Logger,
	tracer tracing.Tracer) cqrs.EventHandle[SendMailEvent] {
	return &sendMailEvent{
		app:    app,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

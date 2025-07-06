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

type deletedEvent struct {
	repo   template.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

func (d *deletedEvent) ToDomain(event template.DefaultEvent) (*template.Domain, error) {
	return d.bs.Template.Existing(template.Params{
		Id:           event.Id,
		Content:      event.Content,
		Subject:      event.Subject,
		Placeholders: event.Placeholders,
		CreatedAt:    event.CreatedAt,
		UpdatedAt:    event.UpdatedAt,
	})
}

// Handle implements cqrs.EventHandle.
func (d *deletedEvent) Handle(ctx context.Context, event template.DefaultEvent) error {
	ctx, span := d.tracer.Start(ctx, "app.template.event.deleted.handler", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", event)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// cast to domain
	payload, err := d.ToDomain(event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to cast to domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return err
	}

	// save to db
	if err := d.repo.Delete(ctx, *payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to delete template",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return err
	}

	d.log.Info("delete process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "DELETE"),
		zap.Any("payload", event),
	)
	return nil
}

func newDeletedEvent(bs domain.Domain, repo template.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.EventHandle[template.DefaultEvent] {
	return &deletedEvent{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

package email

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	"github.com/ose-micro/mailer"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/email"
	"github.com/ose-micro/postman/internal/infrastructure/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type emailApp struct {
	log    logger.Logger
	tracer tracing.Tracer
	create cqrs.CommandHandle[email.CreateCommand, *email.Domain]
	resend cqrs.CommandHandle[email.IdCommand, *email.Domain]
	read   cqrs.QueryHandle[email.ReadQuery, map[string]any]
}

// Create implements email.App.
func (e *emailApp) Create(ctx context.Context, command email.CreateCommand) (*email.Domain, error) {
	ctx, span := e.tracer.Start(ctx, "app.template.create.command", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := e.create.Handle(ctx, command)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	return record, nil
}

// Read implements email.App.
func (e *emailApp) Read(ctx context.Context, request dto.Request) (map[string]any, error) {
	ctx, span := e.tracer.Start(ctx, "app.email.repository.query", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	records, err := e.read.Handle(ctx, email.ReadQuery{
		Request: request,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	return records, nil
}

// Resend implements email.App.
func (r *emailApp) Resend(ctx context.Context, command email.IdCommand) error {
	ctx, span := r.tracer.Start(ctx, "app.email.resend.command", trace.WithAttributes(
		attribute.String("operation", "resend"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := r.resend.Handle(ctx, command); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend"),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func NewEmailApp(bs business.Domain, log logger.Logger, tracer tracing.Tracer,
	read repository.Repository, bus bus.Bus, mailer *mailer.Mailer) email.App {
	return &emailApp{
		log:    log,
		tracer: tracer,
		create: newCreateCommandHandler(bs, read, log, tracer, bus, mailer),
		resend: newResendCommandHandler(bs, read, log, tracer, bus, mailer),
		read:   newReadQueryHandler(read.Email, log, tracer),
	}
}

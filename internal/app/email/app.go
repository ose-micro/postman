package email

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/moriba-cloud/ose-postman/internal/repository/write"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	"github.com/ose-micro/mailer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type emailApp struct {
	log      logger.Logger
	tracer   tracing.Tracer
	create   cqrs.CommandHandle[CreateCommand, email.Domain]
	delete   cqrs.CommandHandle[DeleteCommand, email.Domain]
	resend   cqrs.CommandHandle[ResendCommand, email.Domain]
	read     cqrs.QueryHandle[ReadQuery, []email.Domain]
	read_one cqrs.QueryHandle[ReadOneQuery, email.Domain]
}

// Create implements email.App.
func (e *emailApp) Create(ctx context.Context, command cqrs.Command) (email.Domain, error) {
	ctx, span := e.tracer.Start(ctx, "app.template.create.command", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", command.(CreateCommand))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := e.create.Handle(ctx, command.(CreateCommand))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	return record, nil
}

// Delete implements email.App.
func (e *emailApp) Delete(ctx context.Context, command cqrs.Command) error {
	ctx, span := e.tracer.Start(ctx, "app.template.delete.command", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", command.(CreateCommand))),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := e.create.Handle(ctx, command.(CreateCommand)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return err
	}

	return nil
}

// Read implements email.App.
func (e *emailApp) Read(ctx context.Context, request dto.Request) ([]email.Domain, error) {
	ctx, span := e.tracer.Start(ctx, "app.template.delete.command", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	records, err := e.read.Handle(ctx, ReadQuery{
		request: request,
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

// ReadOne implements email.App.
func (r *emailApp) ReadOne(ctx context.Context, filters ...dto.Filter) (email.Domain, error) {
	return r.read_one.Handle(ctx, ReadOneQuery{
		filters: filters,
	})
}

// Resend implements email.App.
func (r *emailApp) Resend(ctx context.Context, command cqrs.Command) (email.Domain, error) {
	return r.resend.Handle(ctx, command.(ResendCommand))
}

func NewEmailApp(bs domain.Domain, log logger.Logger, tracer tracing.Tracer, write write.Repository,
	read email.Repository, bus bus.Bus, mailer *mailer.Mailer) email.App {
	return &emailApp{
		log:      log,
		tracer:   tracer,
		create:   newCreateCommandHandler(bs, write, log, tracer, bus, mailer),
		delete:   newDeleteCommandHandler(bs, write.Email, log, tracer, bus),
		resend:   newResendCommandHandler(bs, write.Email, log, tracer, bus, mailer),
		read:     newReadQueryHandler(read, log, tracer),
		read_one: newReadOneQueryHandler(read, log, tracer),
	}
}

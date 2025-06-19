package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
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

type ResendCommand struct {
	Id string
}

// CommandName implements cqrs.Command.
func (r ResendCommand) CommandName() string {
	return email.UPDATED_COMMAND
}

// Validate implements cqrs.Command.
func (r ResendCommand) Validate() error {
	fields := make([]string, 0)

	if r.Id == "" {
		fields = append(fields, "id is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}

// Handler
type resendCommandHandler struct {
	repo   email.Repository
	log    logger.Logger
	bus    bus.Bus
	mailer *mailer.Mailer
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (u *resendCommandHandler) Handle(ctx context.Context, command ResendCommand) (email.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.email.resend_mail.command.handler", trace.WithAttributes(
		attribute.String("operation", "RESEND_MAIL"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	var state email.State

	// validate command payload
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "RESEND_MAIL"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	record, err := u.repo.ReadOne(ctx, dto.Filter{
		Field:    "id",
		Operator: dto.EQUAL,
		Value:    command.Id,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to read email",
			zap.String("trace_id", traceId),
			zap.String("operation", "RESEND_MAIL"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	if record.GetState() == email.StateComplete {
		err = fmt.Errorf("email is already in a completed state")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to send mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "RESEND_MAIL"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	err = u.mailer.Send(ctx, mailer.Params{
		Sender:    *record.GetSender(),
		Recipient: record.GetRecipient(),
		Subject:   record.GetSubject(),
		Message:   record.GetMessage(),
		Data:      record.GetData(),
		From:      record.GetFrom(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to resend mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "RESEND_MAIL"),
			zap.Error(err),
		)

		state = email.StateFailed
	} else {
		state = email.StateComplete
	}
	record.SetState(state)

	// save email to write store
	err = u.repo.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to update to postgres",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	// publish bus
	err = u.bus.Publish(command.CommandName(), UpdatedEvent(record.MakePublic()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("publish email updated",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	u.log.Info("update process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "UPDATE"),
		zap.Any("payload", command),
	)
	return *record, nil
}

func newResendCommandHandler(bs domain.Domain, repo email.Repository, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus, mailer *mailer.Mailer) cqrs.CommandHandle[ResendCommand, email.Domain] {
	return &resendCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
		mailer: mailer,
	}
}

package email

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	"github.com/ose-micro/mailer"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/email"
	"github.com/ose-micro/postman/internal/infrastructure/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Handler
type resendCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	bus    domain.Bus
	mailer *mailer.Mailer
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (u *resendCommandHandler) Handle(ctx context.Context, command email.IdCommand) (*email.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.email.resend_mail.command.handler", trace.WithAttributes(
		attribute.String("operation", "resend_mail"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	var state email.State

	// validate command payload
	if err := command.Validate(); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend_mail"),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := u.repo.Email.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Id,
					},
				},
			},
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to repository email",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend_mail"),
			zap.Error(err),
		)

		return nil, err
	}

	if record.State() == email.StateComplete {
		err = ose_error.New(ose_error.ErrNotFound, "email is already in a completed state", traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to send mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend_mail"),
			zap.Error(err),
		)

		return nil, err
	}

	err = u.mailer.Send(ctx, mailer.Params{
		Sender:    record.Sender(),
		Recipient: record.Recipient(),
		Subject:   record.Subject(),
		Message:   record.Message(),
		Data:      record.Data(),
		From:      record.From(),
	})
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to resend mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend_mail"),
			zap.Error(err),
		)

		state = email.StateFailed
	} else {
		state = email.StateComplete
	}
	record.SetState(state)

	// save email to write store
	err = u.repo.Email.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to update to postgres",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend_mail"),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info("update process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "resend_mail"),
		zap.Any("payload", command),
	)
	return record, nil
}

func newResendCommandHandler(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer, bus domain.Bus, mailer *mailer.Mailer) cqrs.CommandHandle[email.IdCommand, *email.Domain] {
	return &resendCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
		mailer: mailer,
	}
}

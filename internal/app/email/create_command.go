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
type createCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	mailer *mailer.Mailer
	bus    domain.Bus
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (c *createCommandHandler) Handle(ctx context.Context, command email.CreateCommand) (*email.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.email.create.command.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process fail",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Any("details", err),
		)

		return nil, err
	}

	temp, err := c.repo.Template.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Template,
					},
				},
			},
		},
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to repository template",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	if err := c.mailer.ValidateMapData(temp.Content(), command.Data); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to validate data",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	state := email.StateFailed
	err = c.mailer.Send(ctx, mailer.Params{
		Sender:    command.Sender,
		Recipient: command.Recipient,
		Subject:   temp.Subject(),
		Message:   temp.Content(),
		Data:      command.Data,
		From:      command.From,
	})

	if err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to send mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		state = email.StateFailed
	} else {
		state = email.StateComplete
	}

	message, err := c.mailer.Rerendered(temp.Content(), command.Data)
	if err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to rendered",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	// create business
	record, err := c.bs.Email.New(email.Params{
		Recipient: command.Recipient,
		Sender:    command.Sender,
		Subject:   temp.Subject(),
		Data:      command.Data,
		Template:  command.Template,
		From:      command.From,
		Message:   message,
		State:     state,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create business",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	// save role to write store
	err = c.repo.Email.Create(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("fail while saving business",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	c.log.Info("create process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("payload", command),
	)
	return record, nil
}

func newCreateCommandHandler(bs business.Domain, repo repository.Repository, log logger.Logger,
	tracer tracing.Tracer, bus domain.Bus, mailer *mailer.Mailer) cqrs.CommandHandle[email.CreateCommand, *email.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
		mailer: mailer,
	}
}

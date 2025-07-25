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

// Handler
type createCommandHandler struct {
	repo   write.Repository
	log    logger.Logger
	mailer *mailer.Mailer
	bus    bus.Bus
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (c *createCommandHandler) Handle(ctx context.Context, command email.CreateCommand) (email.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "app.email.create.command.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("validation process fail",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Any("details", err),
		)

		return email.Domain{}, err
	}

	temp, err := c.repo.Template.Read(ctx, dto.Query{
		Filters: []dto.Filter{
			{
				Field: "id",
				Op: dto.OpEq,
				Value: command.Template,
			},
		},
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to read template",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	if err := c.mailer.ValidateMapData(temp.GetContent(), command.Data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to validate data",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	state := email.StateFailed
	err = c.mailer.Send(ctx, mailer.Params{
		Sender:    command.Sender,
		Recipient: command.Recipient,
		Subject:   temp.GetSubject(),
		Message:   temp.GetContent(),
		Data:      command.Data,
		From:      command.From,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to send mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		state = email.StateFailed
	} else {
		state = email.StateComplete
	}

	message, err := c.mailer.Rerendered(temp.GetContent(), command.Data)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to rendered",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	// create domain
	record, err := c.bs.Email.New(email.Params{
		Recipient: command.Recipient,
		Sender:    command.Sender,
		Subject:   temp.GetSubject(),
		Data:      command.Data,
		Template:  command.Template,
		From:      command.From,
		Message:   message,
		State:     state,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	// save role to write store
	err = c.repo.Email.Create(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("fail while saving domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	// publish bus
	err = c.bus.Publish(command.CommandName(), email.DomainEvent(record.MakePublic()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("publish role created",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return email.Domain{}, err
	}

	c.log.Info("create process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.Any("payload", command),
	)
	return *record, nil
}

func newCreateCommandHandler(bs domain.Domain, repo write.Repository, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus, mailer *mailer.Mailer) cqrs.CommandHandle[email.CreateCommand, email.Domain] {
	return &createCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
		mailer: mailer,
	}
}

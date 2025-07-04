package template

import (
	"context"
	"fmt"
	"strings"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type DeleteCommand struct {
	ID string
}

// CommandName implements cqrs.Command.
func (d DeleteCommand) CommandName() string {
	return template.DELETED_COMMAND
}

// Validate implements cqrs.Command.
func (d DeleteCommand) Validate() error {
	fields := make([]string, 0)

	if d.ID == "" {
		fields = append(fields, "id is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = DeleteCommand{}

// Handler
type deleteCommandHandler struct {
	repo   template.Repository
	log    logger.Logger
	bus    bus.Bus
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (d *deleteCommandHandler) Handle(ctx context.Context, command DeleteCommand) (template.Domain, error) {
	ctx, span := d.tracer.Start(ctx, "app.template.delete.command.handler", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Any("details", err),
		)

		return template.Domain{}, err
	}

	record, err := d.repo.ReadOne(ctx, dto.Filter{
		Field:    "id",
		Operator: dto.EQUAL,
		Value:    command.ID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to read template",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	// save template to write store
	err = d.repo.Delete(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to cast to domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	// publish bus
	err = d.bus.Publish(command.CommandName(), DeletedEvent(record.MakePublic()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("publish template created",
			zap.String("trace_id", traceId),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	d.log.Info("delete process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "DELETE"),
		zap.Any("payload", command),
	)
	return template.Domain{}, nil
}

func newDeleteCommandHandler(bs domain.Domain, repo template.Repository, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus) cqrs.CommandHandle[DeleteCommand, template.Domain] {
	return &deleteCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
	}
}

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

type UpdateCommand struct {
	Id           string
	Content      string
	Subject      string
	Placeholders []string
}

// CommandName implements cqrs.Command.
func (u UpdateCommand) CommandName() string {
	return template.UPDATED_COMMAND
}

// Validate implements cqrs.Command.
func (u UpdateCommand) Validate() error {
	fields := make([]string, 0)

	if u.Id == "" {
		fields = append(fields, "id is required")
	}

	if u.Content == "" {
		fields = append(fields, "content is required")
	}

	if u.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}

// Handler
type updateCommandHandler struct {
	repo   template.Repository
	log    logger.Logger
	bus    bus.Bus
	tracer tracing.Tracer
	bs     domain.Domain
}

// Handle implements cqrs.CommandHandle.
func (u *updateCommandHandler) Handle(ctx context.Context, command UpdateCommand) (template.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "app.template.update.command.handler", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// validate command payload
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process fail",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Any("details", err),
		)

		return template.Domain{}, err
	}

	record, err := u.repo.ReadOne(ctx, dto.Filter{
		Field:    "id",
		Operator: dto.EQUAL,
		Value:    command.Id,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to read template",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Any("details", err),
		)

		return template.Domain{}, err
	}

	if check, _ := u.repo.ReadOne(ctx, dto.Filter{
		Field:    "subject",
		Operator: dto.EQUAL,
		Value:    command.Subject,
	}); check != nil && check.GetID() != record.GetID() {
		err := fmt.Errorf("template already exist with this subject")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to read template",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Any("details", err),
		)

		return template.Domain{}, err
	}

	if err := record.Update(template.Params{
		Content: command.Content,
		Subject: command.Subject,
		Placeholders: command.Placeholders,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to update domain",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Any("details", err),
		)

		return template.Domain{}, err
	}

	// save template to write store
	err = u.repo.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to update to postgres",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	// publish bus
	err = u.bus.Publish(command.CommandName(), UpdatedEvent(record.MakePublic()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("publish template updated",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	u.log.Info("update process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "UPDATE"),
		zap.Any("payload", command),
	)
	return *record, nil
}

func newUpdateCommandHandler(bs domain.Domain, repo template.Repository, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus) cqrs.CommandHandle[UpdateCommand, template.Domain] {
	return &updateCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
	}
}

package template

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	ose_error "github.com/ose-micro/error"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/template"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Handler
type updateCommandHandler struct {
	repo   template.Repo
	log    logger.Logger
	bus    bus.Bus
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (u *updateCommandHandler) Handle(ctx context.Context, command template.UpdateCommand) (bool, error) {
	ctx, span := u.tracer.Start(ctx, "app.template.update.command.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
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
			zap.String("operation", "update"),
			zap.Any("details", err),
		)

		return false, err
	}

	record, err := u.repo.ReadOne(ctx, dto.Request{
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
		u.log.Error("failed to repository template",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Any("details", err),
		)

		return false, err
	}

	if check, _ := u.repo.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Subject,
					},
				},
			},
		},
	}); check != nil && check.ID() != record.ID() {
		err := ose_error.Wrap(err, ose_error.ErrConflict, "template already exist with this subject", traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to repository template",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Any("details", err),
		)

		return false, err
	}

	record.SetContent(command.Content).
		SetSubject(command.Subject).
		SetPlaceholders(command.Placeholders)

	// save template to write store
	err = u.repo.Update(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to update to postgres",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return false, err
	}

	u.log.Info("update process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("payload", command),
	)

	return true, nil
}

func newUpdateCommandHandler(bs business.Domain, repo template.Repo, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus) cqrs.CommandHandle[template.UpdateCommand, bool] {
	return &updateCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
	}
}

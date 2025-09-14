package template

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/cqrs/bus"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/template"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Handler
type deleteCommandHandler struct {
	repo   template.Repo
	log    logger.Logger
	bus    bus.Bus
	tracer tracing.Tracer
	bs     business.Domain
}

// Handle implements cqrs.CommandHandle.
func (d *deleteCommandHandler) Handle(ctx context.Context, command template.DeleteCommand) (bool, error) {
	ctx, span := d.tracer.Start(ctx, "app.template.delete.command.handler", trace.WithAttributes(
		attribute.String("operation", "delete"),
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
			zap.String("operation", "delete"),
			zap.Any("details", err),
		)

		return false, err
	}

	record, err := d.repo.ReadOne(ctx, dto.Request{
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
		d.log.Error("failed to repository template",
			zap.String("trace_id", traceId),
			zap.String("operation", "delete"),
			zap.Error(err),
		)

		return false, err
	}

	// save template to write store
	err = d.repo.Delete(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		d.log.Error("failed to cast to business",
			zap.String("trace_id", traceId),
			zap.String("operation", "delete"),
			zap.Error(err),
		)

		return false, err
	}

	d.log.Info("delete process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "delete"),
		zap.Any("payload", command),
	)

	return true, nil
}

func newDeleteCommandHandler(bs business.Domain, repo template.Repo, log logger.Logger,
	tracer tracing.Tracer, bus bus.Bus) cqrs.CommandHandle[template.DeleteCommand, bool] {
	return &deleteCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bus:    bus,
		bs:     bs,
	}
}

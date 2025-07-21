package template

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type readQueryHandler struct {
	repo   template.Read
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readQueryHandler) Handle(ctx context.Context, query template.ReadQuery) (map[string]any, error) {
	ctx, span := r.tracer.Start(ctx, "app.template.read.query.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch roles from store
	records, err := r.repo.Read(ctx, query.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read roles",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)

		return nil, err
	}

	r.log.Info("read process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "READ"),
		zap.Any("payload", fmt.Sprintf("%v", query)),
	)
	return records, nil
}

func newReadQueryHandler(repo template.Read, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[template.ReadQuery, map[string]any] {
	return &readQueryHandler{
		repo: repo,
		log: log,
		tracer: tracer,
	}
}

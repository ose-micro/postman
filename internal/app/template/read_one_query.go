package template

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ReadOneQuery struct {
	filters []dto.Filter
}

// QueryName implements cqrs.Query.
func (r ReadOneQuery) QueryName() string {
	return template.READ_QUERY
}

var _ cqrs.Query = ReadQuery{}

// Handler
type readOneQueryHandler struct {
	repo   template.Repository
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readOneQueryHandler) Handle(ctx context.Context, query ReadOneQuery) (template.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "app.template.read_one.query.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch templates from store
	record, err := r.repo.ReadOne(ctx, query.filters...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read template",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ_ONE"),
			zap.Error(err),
		)

		return template.Domain{}, err
	}

	r.log.Info("read_one process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "READ_ONE"),
		zap.Any("payload", fmt.Sprintf("%v", query)),
	)
	return *record, nil
}

func newReadOneQueryHandler(repo template.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[ReadOneQuery, template.Domain] {
	return &readOneQueryHandler{
		repo: repo,
		log: log,
		tracer: tracer,
	}
}

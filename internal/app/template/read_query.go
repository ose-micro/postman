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

type ReadQuery struct {
	request dto.Request
}

// QueryName implements cqrs.Query.
func (r ReadQuery) QueryName() string {
	return template.READ_QUERY
}

var _ cqrs.Query = ReadQuery{}

// Handler
type readQueryHandler struct {
	repo   template.Repository
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readQueryHandler) Handle(ctx context.Context, query ReadQuery) ([]template.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "app.template.read.query.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch templates from store
	records, err := r.repo.Read(ctx, query.request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read templates",
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

func newReadQueryHandler(repo template.Repository, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[ReadQuery, []template.Domain] {
	return &readQueryHandler{
		repo: repo,
		log: log,
		tracer: tracer,
	}
}

package email

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"github.com/ose-micro/postman/internal/business/email"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type readQueryHandler struct {
	repo   email.Repo
	log    logger.Logger
	tracer tracing.Tracer
}

// Handle implements cqrs.QueryHandle.
func (r *readQueryHandler) Handle(ctx context.Context, query email.ReadQuery) (map[string]any, error) {
	ctx, span := r.tracer.Start(ctx, "app.email.repository.query.handler", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("payload", fmt.Sprintf("%v", query)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// fetch roles from store
	records, err := r.repo.Read(ctx, query.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository roles",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)

		return nil, err
	}

	r.log.Info("repository process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "read"),
		zap.Any("payload", fmt.Sprintf("%v", query)),
	)
	return records, nil
}

func newReadQueryHandler(repo email.Repo, log logger.Logger,
	tracer tracing.Tracer) cqrs.QueryHandle[email.ReadQuery, map[string]any] {
	return &readQueryHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
	}
}

package template

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/postgres"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type templateRepository struct {
	db     *postgres.Postgres
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Create implements template.Repository.
func (c *templateRepository) Create(ctx context.Context, payload template.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.write.template.create", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	entity := newEntity(payload)
	if err := c.db.Conn().Save(entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("create process complete successfully",
		zap.String("operation", "CREATE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (c *templateRepository) Delete(ctx context.Context, payload template.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.write.template.delete", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	model := newEntity(payload)
	if err := c.db.Conn().Delete(model).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to delete in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("delete process complete successfully",
		zap.String("operation", "DELETE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}
func (c *templateRepository) Read(ctx context.Context, request dto.Query) (*template.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "repository.write.template.read", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	var entity Template
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	query := c.db.Conn()
	query, err := postgres.BuildGORMQuery(query, "templates", request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to build query",
			zap.String("trace_id", traceID),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	if err := query.First(&entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to read templates",
			zap.String("operation", "READ"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return nil, err
	}

	c.log.Info("read process complete successfully",
		zap.String("operation", "READ"),
		zap.String("trace_id", traceID),
		zap.Any("payload", request),
	)

	return c.toDomain(entity)
}

func (c *templateRepository) Update(ctx context.Context, payload template.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.write.template.update", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	entity := newEntity(payload)
	if err := c.db.Conn().Save(entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to update in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("update process complete successfully",
		zap.String("operation", "UPDATE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (c *templateRepository) toDomain(entity Template) (*template.Domain, error) {
	return c.bs.Template.Existing(template.Params{
		Id:           entity.ID,
		Content:      entity.Content,
		Subject:      entity.Subject,
		Placeholders: entity.Placeholders,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	})
}

func NewRepository(db *postgres.Postgres, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) template.Write {
	return &templateRepository{
		db:     db,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

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
func (t *templateRepository) Create(ctx context.Context, payload template.Domain) error {
	ctx, span := t.tracer.Start(ctx, "repository.write.template.create", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	entity := newEntity(payload)
	if err := t.db.Conn().Save(entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to create in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	t.log.Info("create process complete successfully",
		zap.String("operation", "CREATE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (t *templateRepository) Delete(ctx context.Context, payload template.Domain) error {
	ctx, span := t.tracer.Start(ctx, "repository.write.template.delete", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	model := newEntity(payload)
	if err := t.db.Conn().Delete(model).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to delete in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)
		return err
	}

	t.log.Info("delete process complete successfully",
		zap.String("operation", "DELETE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (t *templateRepository) Read(ctx context.Context, request dto.Request) ([]template.Domain, error) {
	ctx, span := t.tracer.Start(ctx, "repository.write.template.read", trace.WithAttributes(
		attribute.String("action", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	query := t.db.Conn()
	query = postgres.BuildFilters[Template](query, request.Filter...)
	query = postgres.BuildSort(query, request.Sort...)

	if request.Pagination != nil {
		if request.Pagination.Page > 0 {
			offset := (request.Pagination.Page - 1) * request.Pagination.Limit
			query = query.Limit(request.Pagination.Limit).Offset(offset)
		}
	}

	entities := make([]Template, 0)
	if err := query.Find(&entities).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to read templates",
			zap.String("trace_id", traceID),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	records := make([]template.Domain, len(entities))

	for i, record := range entities {
		v, err := t.toDomain(record)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			t.log.Error("failed to cast to domain",
				zap.String("trace_id", traceID),
				zap.String("operation", "READ"),
				zap.Error(err),
			)

			return nil, err
		}

		records[i] = *v
	}

	t.log.Info("read process create successfully",
		zap.String("operation", "READ"),
		zap.String("trace_id", traceID),
		zap.Any("payload", request),
	)

	return records, nil
}

func (t *templateRepository) ReadOne(ctx context.Context, filters ...dto.Filter) (*template.Domain, error) {
	ctx, span := t.tracer.Start(ctx, "repository.write.template.read_one", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", filters)),
	))
	defer span.End()

	var entity Template
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	query := t.db.Conn()
	query = postgres.BuildFilters[Template](query, filters...)

	if err := query.First(&entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to read template",
			zap.String("operation", "READ_ONE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return nil, err
	}

	t.log.Info("read_one process complete successfully",
		zap.String("operation", "READ_ONE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", filters),
	)

	return t.toDomain(entity)
}

func (r *templateRepository) Update(ctx context.Context, payload template.Domain) error {
	ctx, span := r.tracer.Start(ctx, "repository.write.template.update", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	entity := newEntity(payload)
	if err := r.db.Conn().Save(entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to update in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)
		return err
	}

	r.log.Info("update process complete successfully",
		zap.String("operation", "UPDATE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (t *templateRepository) toDomain(entity Template) (*template.Domain, error) {
	return t.bs.Template.Existing(template.Params{
		Id:           entity.ID,
		Content:      entity.Content,
		Subject:      entity.Subject,
		Placeholders: entity.Placeholders,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	})
}

func NewTemplateRepository(db *postgres.Postgres, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) template.Repository {
	return &templateRepository{
		db:     db,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

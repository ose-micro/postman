package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/postgres"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type emailRepository struct {
	db     *postgres.Postgres
	log    logger.Logger
	tracer tracing.Tracer
	bs     domain.Domain
}

// Create implements role.Repository.
func (r *emailRepository) Create(ctx context.Context, payload email.Domain) error {
	ctx, span := r.tracer.Start(ctx, "repository.write.email.create", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	entity := newEntity(payload)
	if err := r.db.Conn().Save(entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to create in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	r.log.Info("create process complete successfully",
		zap.String("operation", "CREATE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (r *emailRepository) Delete(ctx context.Context, payload email.Domain) error {
	ctx, span := r.tracer.Start(ctx, "repository.write.email.delete", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	model := newEntity(payload)
	if err := r.db.Conn().Delete(model).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to delete in postgres",
			zap.String("trace_id", traceID),
			zap.String("operation", "DELETE"),
			zap.Error(err),
		)
		return err
	}

	r.log.Info("delete process complete successfully",
		zap.String("operation", "DELETE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

func (r *emailRepository) Read(ctx context.Context, request dto.Request) ([]email.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "repository.write.role.read", trace.WithAttributes(
		attribute.String("action", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	query := r.db.Conn()
	query = postgres.BuildFilters[Email](query, request.Filter...)
	query = postgres.BuildSort(query, request.Sort...)

	if request.Pagination != nil {
		if request.Pagination.Page > 0 {
			offset := (request.Pagination.Page - 1) * request.Pagination.Limit
			query = query.Limit(request.Pagination.Limit).Offset(offset)
		}
	}

	entities := make([]Email, 0)
	if err := query.Find(&entities).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read emails",
			zap.String("trace_id", traceID),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	records := make([]email.Domain, len(entities))

	for i, record := range entities {
		v, err := r.toDomain(record, traceID, span)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			r.log.Error("failed to cast to domain",
				zap.String("trace_id", traceID),
				zap.String("operation", "READ"),
				zap.Error(err),
			)

			return nil, err
		}

		records[i] = *v
	}

	r.log.Info("read process create successfully",
		zap.String("operation", "READ"),
		zap.String("trace_id", traceID),
		zap.Any("payload", request),
	)

	return records, nil
}

func (r *emailRepository) ReadOne(ctx context.Context, filters ...dto.Filter) (*email.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "repository.write.email.read_one", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", filters)),
	))
	defer span.End()

	var entity Email
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	query := r.db.Conn()
	query = postgres.BuildFilters[Email](query, filters...)

	if err := query.First(&entity).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to read email",
			zap.String("operation", "READ_ONE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return nil, err
	}

	r.log.Info("read_one process complete successfully",
		zap.String("operation", "READ_ONE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", filters),
	)

	return r.toDomain(entity, traceID, span)
}

func (r *emailRepository) Update(ctx context.Context, payload email.Domain) error {
	ctx, span := r.tracer.Start(ctx, "repository.write.email.update", trace.WithAttributes(
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

func (r *emailRepository) toDomain(entity Email, traceId string, span trace.Span) (*email.Domain, error) {
	return r.bs.Email.Existing(email.Params{
		Id:        entity.Id,
		Recipient: entity.Recipient,
		Sender:    entity.Sender,
		Subject:   entity.Subject,
		Data: func() map[string]interface{} {
			var data map[string]interface{}
			err := json.Unmarshal([]byte(entity.Data), &data)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				r.log.Fatal("failed to cast to domain",
					zap.String("trace_id", traceId),
					zap.Error(err))
			}

			return data
		}(),
		Template:  entity.Template,
		From:      entity.From,
		Message:   entity.Message,
		State:     entity.State,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	})
}

func NewEmailRepository(db *postgres.Postgres, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) email.Repository {
	return &emailRepository{
		db:     db,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

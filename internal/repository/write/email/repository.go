package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

// Create implements email.Repository.
func (c *emailRepository) Create(ctx context.Context, payload email.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.write.email.create", trace.WithAttributes(
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

func (c *emailRepository) Delete(ctx context.Context, payload email.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.write.email.delete", trace.WithAttributes(
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
func (c *emailRepository) Read(ctx context.Context, request dto.Query) (*email.Domain, error) {
	ctx, span := c.tracer.Start(ctx, "repository.write.email.read", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	var entity Email
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	query := c.db.Conn()
	query, err := postgres.BuildGORMQuery(query, "emails", request)
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
		c.log.Error("failed to read emails",
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

func (c *emailRepository) Update(ctx context.Context, payload email.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.write.email.update", trace.WithAttributes(
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

func (c *emailRepository) toDomain(entity Email) (*email.Domain, error) {
	return c.bs.Email.Existing(email.Params{
		Id:             entity.Id,
		Recipient: entity.Recipient,
		Sender: entity.Sender,
		Subject: entity.Subject,
		Data: func() map[string]interface{} {
			if entity.Data == "" {
				return nil
			}

			var value map[string]interface{}
			if err := json.Unmarshal([]byte(entity.Data), &value); err != nil {
				log.Fatalln(err)
			}

			return value
		}(),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	})
}

func NewRepository(db *postgres.Postgres, bs domain.Domain, log logger.Logger, tracer tracing.Tracer) email.Write {
	return &emailRepository{
		db:     db,
		log:    log,
		tracer: tracer,
		bs:     bs,
	}
}

package email

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	mongodb "github.com/ose-micro/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type emailRepository struct {
	collection *mongo.Collection
	log        logger.Logger
	tracer     tracing.Tracer
	bs         domain.Domain
}

// Create implements email.Repository.
func (c *emailRepository) Create(ctx context.Context, payload email.Domain) error {
	ctx, span := c.tracer.Start(ctx, "read.repository.email.create", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	record := newCollection(payload)
	if _, err := c.collection.InsertOne(ctx, record); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create in mongo",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("create process complete successfully",
		zap.String("operation", "CREATE"),
		zap.String("trace_id", traceId),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

// Delete implements email.Repository.
func (c *emailRepository) Delete(ctx context.Context, payload email.Domain) error {
	ctx, span := c.tracer.Start(ctx, "read.repository.email.delete", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%+v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	filter := bson.M{"_id": payload.GetID()}
	if _, err := c.collection.DeleteOne(ctx, filter); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.log.Error("failed to delete in mongo",
			zap.String("operation", "DELETE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("delete process completed successfully",
		zap.String("operation", "DELETE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)

	return nil
}

// Read implements email.Repository.
func (c *emailRepository) Read(ctx context.Context, request dto.Request) (map[string]any, error) {
	ctx, span := c.tracer.Start(ctx, "read.repository.email.read", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%+v", request)),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	mongodb.RegisterType("email", email.Public{})
	typeHints := map[string]string{}

	for _, v := range request.Queries {
		typeHints[v.Name] = "email"
	}

	res, err := mongodb.RunFaceted(ctx, c.collection, request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.log.Error("Failed to fetch role by request",
			zap.String("operation", "READ"),
			zap.String("trace_id", traceID),
			zap.Any("payload", request),
			zap.Error(err),
		)
		return nil, err
	}

	c.log.Info("Read process completed successfully",
		zap.String("operation", "READ"),
		zap.String("trace_id", traceID),
		zap.Any("payload", request),
	)

	records, err := mongodb.CastFacetedResult(res, typeHints)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("Failed to cast faceted result",
			zap.String("operation", "READ"),
			zap.String("trace_id", traceID),
			zap.Any("payload", request),
			zap.Error(err),
		)
		return nil, err
	}

	return records, nil
}

// Update implements email.Repository.
func (c *emailRepository) Update(ctx context.Context, payload email.Domain) error {
	ctx, span := c.tracer.Start(ctx, "repository.read.email.update", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%+v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	collection := newCollection(payload)
	filter := bson.M{"_id": payload.GetID()}

	if _, err := c.collection.UpdateOne(ctx, filter, bson.M{
		"$set": collection,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.log.Error("failed to update email",
			zap.String("operation", "UPDATE"),
			zap.String("trace_id", traceID),
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

func NewRepository(db *mongodb.Client, log logger.Logger, tracer tracing.Tracer, bs domain.Domain) email.Read {
	return &emailRepository{
		log:        log,
		tracer:     tracer,
		bs:         bs,
		collection: db.Collection("emails"),
	}
}
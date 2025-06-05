package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	mongodb "github.com/ose-micro/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
func (e *emailRepository) Create(ctx context.Context, payload email.Domain) error {
	ctx, span := e.tracer.Start(ctx, "read.repository.email.create", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", payload.MakePublic())),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	record := newCollection(payload)
	if _, err := e.collection.InsertOne(ctx, record); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to create in mongo",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
		return err
	}

	e.log.Info("create process complete successfully",
		zap.String("operation", "CREATE"),
		zap.String("trace_id", traceId),
		zap.Any("payload", payload.MakePublic()),
	)
	return nil
}

// Delete implements email.Repository.
func (e *emailRepository) Delete(ctx context.Context, payload email.Domain) error {
	ctx, span := e.tracer.Start(ctx, "read.repository.email.delete", trace.WithAttributes(
		attribute.String("operation", "DELETE"),
		attribute.String("payload", fmt.Sprintf("%+v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	filter := bson.M{"_id": payload.GetID()}
	if _, err := e.collection.DeleteOne(ctx, filter); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		e.log.Error("failed to delete in mongo",
			zap.String("operation", "DELETE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return err
	}

	e.log.Info("delete process completed successfully",
		zap.String("operation", "DELETE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)

	return nil
}

// Read implements email.Repository.
func (e *emailRepository) Read(ctx context.Context, request dto.Request) ([]email.Domain, error) {
	ctx, span := e.tracer.Start(ctx, "read.repository.email.read", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%+v", request)),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	filters := mongodb.BuildFilter(request.Filter)

	findOpts := options.Find()

	if len(request.Sort) > 0 {
		sorts := make([]mongodb.Sort, 0)
		for _, s := range request.Sort {
			sorts = append(sorts, mongodb.Sort{
				Field:     s.Field,
				Direction: mongodb.Direction(s.Value),
			})
		}

		sort := mongodb.BuildSort(sorts...)
		sort(findOpts)
	}

	if request.Pagination != nil {
		limitValue := int64(request.Pagination.Limit)
		limit := mongodb.WithLimit(limitValue)

		skipValue := int64(request.Pagination.Page)
		skip := mongodb.WithSkip(skipValue)

		limit(findOpts)
		skip(findOpts)
	}

	cursor, err := e.collection.Find(ctx, filters, findOpts)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("Failed to fetch email by request",
			zap.String("operation", "READ"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return nil, err
	}

	defer cursor.Close(ctx)

	var records []email.Domain
	for cursor.Next(ctx) {
		var collection Email
		if err := cursor.Decode(&collection); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			e.log.Error("failed to cast to collection",
				zap.String("operation", "READ"),
				zap.String("trace_id", traceID),
				zap.Error(err),
			)
			return nil, err
		}

		record, err := e.toDomain(collection)
		if err != nil {
			span.RecordError(err)
			e.log.Error("failed to cast to domain",
				zap.String("trace_id", traceID),
				zap.String("operation", "READ"),
				zap.Error(err),
			)
			return nil, err
		}
		records = append(records, *record)
	}

	e.log.Info("read process complete successfully",
		zap.String("operation", "READ"),
		zap.String("trace_id", traceID),
		zap.Any("payload", fmt.Sprintf("%+v", request)),
	)
	return records, nil
}

// ReadOne implements email.Repository.
func (e *emailRepository) ReadOne(ctx context.Context, filters ...dto.Filter) (*email.Domain, error) {
	ctx, span := e.tracer.Start(ctx, "repository.read.email.read_one", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%+v", filters)),
	))
	defer span.End()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	filter := mongodb.BuildFilter(filters)

	var collection Email

	err := e.collection.FindOne(ctx, filter).
		Decode(&collection)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to find email",
			zap.String("operation", "READ_ONE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("no document with this given filters")
	} else if err != nil {
		return nil, err
	}

	record, err := e.toDomain(collection)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to cast to domain",
			zap.String("operation", "READ_ONE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
	}

	e.log.Info("read_one process complete successfully",
		zap.String("operation", "READ_ONE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", filters),
	)
	return record, err
}

// Update implements email.Repository.
func (e *emailRepository) Update(ctx context.Context, payload email.Domain) error {
	ctx, span := e.tracer.Start(ctx, "repository.read.email.update", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%+v", payload.MakePublic())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	collection := newCollection(payload)
	filter := bson.M{"_id": payload.GetID()}

	if _, err := e.collection.UpdateOne(ctx, filter, bson.M{
		"$set": collection,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		e.log.Error("failed to update email",
			zap.String("operation", "UPDATE"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return err
	}

	e.log.Info("update process complete successfully",
		zap.String("operation", "UPDATE"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.MakePublic()),
	)

	return nil
}

func (e *emailRepository) toDomain(collection Email) (*email.Domain, error) {
	return e.bs.Email.Existing(email.Params{
		Id:        collection.Id,
		Recipient: collection.Recipient,
		Sender:    collection.Sender,
		Subject:   collection.Subject,
		Data: func() map[string]interface{} {
			if data, ok := collection.Data.(map[string]interface{}); ok {
				return data
			}
			return nil
		}(),
		Template:  collection.Template,
		From:      collection.From,
		Message:   collection.Message,
		State:     collection.State,
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	})
}

func NewEmailRepository(bs domain.Domain, db *mongodb.Client, log logger.Logger, tracer tracing.Tracer) email.Repository {
	return &emailRepository{
		log:        log,
		tracer:     tracer,
		bs:         bs,
		collection: db.Collection("emails"),
	}
}

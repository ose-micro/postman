package template

import (
	"context"
	"fmt"

	"github.com/ose-micro/common"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_error "github.com/ose-micro/error"
	mongodb "github.com/ose-micro/mongo"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/business/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repository struct {
	collection *mongo.Collection
	log        logger.Logger
	tracer     tracing.Tracer
	bs         business.Domain
}

func (r *repository) ReadOne(ctx context.Context, request dto.Request) (*template.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "infrastructure.repository.template.read_one", trace.WithAttributes(
		attribute.String("operation", "read_one"),
		attribute.String("dto", fmt.Sprintf("%v", request))),
	)
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	res, err := r.Read(ctx, request)
	if err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository res",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)
		return nil, err
	}

	raw, ok := res["one"]
	if !ok {
		return nil, ose_error.New(ose_error.ErrNotFound, "email not found")
	}

	var records []template.Public

	if err := common.JsonToAny(raw, &records); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository res",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)
		return nil, err
	}

	if len(records) == 0 {
		err := ose_error.New(ose_error.ErrNotFound, "template not found", traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository res",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)
		return nil, err
	}

	return r.toDomain(records[0]), nil
}

// Create implements template.Repository.
func (c *repository) Create(ctx context.Context, payload template.Domain) error {
	ctx, span := c.tracer.Start(ctx, "infrastructure.repository.template.create", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", payload.Public())),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	record := newCollection(payload)
	if _, err := c.collection.InsertOne(ctx, record); err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create in mongo",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("create process complete successfully",
		zap.String("operation", "create"),
		zap.String("trace_id", traceId),
		zap.Any("payload", payload.Public()),
	)
	return nil
}

// Delete implements template.Repository.
func (c *repository) Delete(ctx context.Context, payload template.Domain) error {
	ctx, span := c.tracer.Start(ctx, "infrastructure.repository.template.delete", trace.WithAttributes(
		attribute.String("operation", "delete"),
		attribute.String("payload", fmt.Sprintf("%+v", payload.Public())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	filter := bson.M{"_id": payload.ID()}
	if _, err := c.collection.DeleteOne(ctx, filter); err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceID)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to delete in mongo",
			zap.String("operation", "delete"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("delete process completed successfully",
		zap.String("operation", "delete"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.Public()),
	)

	return nil
}

// Read implements template.Repository.
func (c *repository) Read(ctx context.Context, request dto.Request) (map[string]any, error) {
	ctx, span := c.tracer.Start(ctx, "infrastructure.repository.template.repository", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("payload", fmt.Sprintf("%+v", request)),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	mongodb.RegisterType("template", template.Public{})
	typeHints := map[string]string{}

	for _, v := range request.Queries {
		typeHints[v.Name] = "template"
	}

	res, err := mongodb.RunFaceted(ctx, c.collection, request)
	if err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceID)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.log.Error("Failed to fetch role by request",
			zap.String("operation", "read"),
			zap.String("trace_id", traceID),
			zap.Any("payload", request),
			zap.Error(err),
		)
		return nil, err
	}

	c.log.Info("Read process completed successfully",
		zap.String("operation", "read"),
		zap.String("trace_id", traceID),
		zap.Any("payload", request),
	)

	records, err := mongodb.CastFacetedResult(res, typeHints)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("Failed to cast faceted result",
			zap.String("operation", "read"),
			zap.String("trace_id", traceID),
			zap.Any("payload", request),
			zap.Error(err),
		)
		return nil, err
	}

	return records, nil
}

// Update implements template.Repository.
func (c *repository) Update(ctx context.Context, payload template.Domain) error {
	ctx, span := c.tracer.Start(ctx, "infrastructure.repository.template.update", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("payload", fmt.Sprintf("%+v", payload.Public())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	collection := newCollection(payload)
	filter := bson.M{"_id": payload.ID()}

	if _, err := c.collection.UpdateOne(ctx, filter, bson.M{
		"$set": collection,
	}); err != nil {
		err = ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceID)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.log.Error("failed to update template",
			zap.String("operation", "update"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return err
	}

	c.log.Info("update process complete successfully",
		zap.String("operation", "update"),
		zap.String("trace_id", traceID),
		zap.Any("payload", payload.Public()),
	)

	return nil
}

func (r *repository) toDomain(payload template.Public) *template.Domain {
	result, _ := r.bs.Template.Existing(*payload.Params())
	return result
}

func NewRepository(db *mongodb.Client, log logger.Logger, tracer tracing.Tracer, bs business.Domain) template.Repo {
	return &repository{
		log:        log,
		tracer:     tracer,
		bs:         bs,
		collection: db.Collection("templates"),
	}
}

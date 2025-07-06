package handlers

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/app"
	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	templateDomain "github.com/moriba-cloud/ose-postman/internal/domain/template"
	templatev1 "github.com/moriba-cloud/ose-postman/internal/interface/grpc/gen/go/template/v1"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	templateHandler struct {
		templatev1.UnimplementedTemplateServiceServer
		app    templateDomain.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (r *templateHandler) response(param templateDomain.Domain) *templatev1.Template {
	return &templatev1.Template{
		Id:        param.GetID(),
		Subject: param.GetSubject(),
		Content: param.GetContent(),
		Placeholders: param.GetPlaceholders(),
		CreatedAt: timestamppb.New(param.GetCreatedAt()),
		UpdatedAt: timestamppb.New(param.GetUpdatedAt()),
	}
}

func (r *templateHandler) Create(ctx context.Context, request *templatev1.CreateRequest) (*templatev1.CreateResponse, error) {
	ctx, span := r.tracer.Start(ctx, "interface.grpc.template.create.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := template.CreateCommand{
		Content:      request.Content,
		Subject:      request.Subject,
		Placeholders: request.Placeholders,
	}

	record, err := r.app.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to create template",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}

	r.log.Info("template create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.Any("payload", request),
	)

	return &templatev1.CreateResponse{
		Message: "template create successfully",
		Record:  r.response(record),
	}, nil
}

func (t *templateHandler) Read(ctx context.Context, request *templatev1.ReadRequest) (*templatev1.ReadResponse, error) {
	ctx, span := t.tracer.Start(ctx, "interface.grpc.template.read.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	query, err := buildAppRequest(request.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to case to dto",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	result, err := t.app.Read(ctx, *query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to read templates",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	records := make([]*templatev1.Template, len(result))
	for i, o := range result {
		records[i] = t.response(o)
	}

	resQuery, err := buildGRPCRequest(&dto.Request{
		Pagination: &dto.Pagination{
			Page:  query.Pagination.Page,
			Limit: query.Pagination.Limit,
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("failed to case to grpc",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	return &templatev1.ReadResponse{
		Records: records,
		Query:   resQuery,
	}, nil
}

func (r *templateHandler) Update(ctx context.Context, request *templatev1.UpdateRequest) (*templatev1.UpdateResponse, error) {
	ctx, span := r.tracer.Start(ctx, "interface.grpc.template.update.handler", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := template.UpdateCommand{
		Id:           request.Id,
		Content:      request.Content,
		Subject:      request.Subject,
		Placeholders: request.Placeholders,
	}

	err := r.app.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to case to grpc",
			zap.String("trace_id", traceId),
			zap.String("operation", "UPDATE"),
			zap.Error(err),
		)
		return nil, err
	}

	return &templatev1.UpdateResponse{
		Message: "template update successfully",
	}, nil
}

func (r *templateHandler) ReadOne(ctx context.Context, request *templatev1.ReadOneRequest) (*templatev1.ReadOneResponse, error) {
	ctx, span := r.tracer.Start(ctx, "interface.grpc.template.read_one.handler", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	filters := make([]dto.Filter, 0)

	for _, filter := range request.Filter {
		filters = append(filters, dto.Filter{
			Field:    filter.Field,
			Operator: enumToOperator(filter.Operator),
			Value:    processValue(filter),
		})
	}

	campaign, err := r.app.ReadOne(ctx, filters...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to case to grpc",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ_ONE"),
			zap.Error(err),
		)
		return nil, err
	}

	return &templatev1.ReadOneResponse{
		Message: "template read successfully",
		Record:  r.response(campaign),
	}, nil
}

func (r *templateHandler) Delete(ctx context.Context, request *templatev1.DeleteRequest) (*templatev1.DeleteResponse, error) {
	ctx, span := r.tracer.Start(ctx, "interface.grpc.template.read_one.handler", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := template.DeleteCommand{
		Id: request.Id,
	}

	err := r.app.Delete(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to case to grpc",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ_ONE"),
			zap.Error(err),
		)
		return nil, err
	}

	return &templatev1.DeleteResponse{
		Message: "template deleted successfully",
	}, nil
}

func NewTemplate(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *templateHandler {
	return &templateHandler{
		app:    apps.Template,
		log:    log,
		tracer: tracer,
	}
}

package handlers

import (
	"context"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/app"
	"github.com/moriba-cloud/ose-postman/internal/app/email"
	emailDomain "github.com/moriba-cloud/ose-postman/internal/domain/email"
	commonv1 "github.com/moriba-cloud/ose-postman/internal/interface/grpc/gen/go/common/v1"
	emailv1 "github.com/moriba-cloud/ose-postman/internal/interface/grpc/gen/go/email/v1"
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
	emailHandler struct {
		emailv1.UnimplementedEmailServiceServer
		app    emailDomain.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (e *emailHandler) response(param emailDomain.Domain) *emailv1.Email {
	return &emailv1.Email{
		Id:        param.GetID(),
		Recipient: param.GetRecipient(),
		Data: func() []*commonv1.Field {
			data := make([]*commonv1.Field, 0)
			for k, v := range param.GetData() {
				data = append(data, &commonv1.Field{
					Field: k,
					Value: v.(string),
				})
			}
			return data
		}(),
		Subject:  param.GetSubject(),
		Sender:   func () string {
			if param.GetSender() == nil {
				return ""
			}

			return param.GetSubject()
		}(),
		From:     param.GetFrom(),
		Template: param.GetTemplate(),
		Message:  param.GetMessage(),
		State: func() emailv1.State {
			switch param.GetState() {
			case emailDomain.StateFailed:
				return emailv1.State_StateFailed
			case emailDomain.StateComplete:
				return emailv1.State_StateComplete
			default:
				return emailv1.State_StateUnknown
			}
		}(),
		CreatedAt: timestamppb.New(param.GetCreatedAt()),
		UpdatedAt: timestamppb.New(param.GetUpdatedAt()),
	}
}

func (e *emailHandler) Create(ctx context.Context, request *emailv1.CreateRequest) (*emailv1.CreateResponse, error) {
	ctx, span := e.tracer.Start(ctx, "interface.grpc.email.create.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := email.CreateCommand{
		Recipient: request.Recipient,
		Sender:    &request.Sender,
		Data: func() map[string]interface{} {
			data := make(map[string]interface{})
			for _, v := range request.Data {
				data[v.Field] = v.Value
			}

			return data
		}(),
		Template: request.Template,
		From:     request.From,
	}

	record, err := e.app.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to create email",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}

	e.log.Info("email create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.Any("payload", request),
	)

	return &emailv1.CreateResponse{
		Message: "email create successfully",
		Record:  e.response(record),
	}, nil
}

func (e *emailHandler) Resend(ctx context.Context, request *emailv1.ResendRequest) (*emailv1.ResendResponse, error) {
	ctx, span := e.tracer.Start(ctx, "interface.grpc.email.resend.handler", trace.WithAttributes(
		attribute.String("operation", "RESEND"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := email.ResendCommand{Id: request.Id}

	if _, err := e.app.Resend(ctx, payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to resend verification mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "RESEND"),
			zap.Error(err),
		)

		return nil, err
	}

	e.log.Info("resend verification mail process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "RESEND"),
		zap.Any("payload", request),
	)

	return &emailv1.ResendResponse{
		Message: "resend mail successfully",
	}, nil
}

func (e *emailHandler) Read(ctx context.Context, request *emailv1.ReadRequest) (*emailv1.ReadResponse, error) {
	ctx, span := e.tracer.Start(ctx, "interface.grpc.email.read.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	query, err := buildAppRequest(request.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to case to dto",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	result, err := e.app.Read(ctx, *query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to read emails",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	records := make([]*emailv1.Email, len(result))
	for i, o := range result {
		records[i] = e.response(o)
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
		e.log.Error("failed to case to grpc",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	return &emailv1.ReadResponse{
		Records: records,
		Request: resQuery,
	}, nil
}

func (r *emailHandler) ReadOne(ctx context.Context, request *emailv1.ReadOneRequest) (*emailv1.ReadOneResponse, error) {
	ctx, span := r.tracer.Start(ctx, "interface.grpc.email.read_one.handler", trace.WithAttributes(
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

	return &emailv1.ReadOneResponse{
		Message: "email read successfully",
		Record:  r.response(campaign),
	}, nil
}

func (r *emailHandler) Delete(ctx context.Context, request *emailv1.DeleteRequest) (*emailv1.DeleteResponse, error) {
	ctx, span := r.tracer.Start(ctx, "interface.grpc.read_one.handler", trace.WithAttributes(
		attribute.String("operation", "READ_ONE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := email.DeleteCommand{
		ID: request.Id,
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

	return &emailv1.DeleteResponse{
		Message: "email deleted successfully",
	}, nil
}

func NewEmail(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *emailHandler {
	return &emailHandler{
		app:    apps.Email,
		log:    log,
		tracer: tracer,
	}
}

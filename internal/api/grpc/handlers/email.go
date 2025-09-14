package handlers

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_error "github.com/ose-micro/error"
	emailv1 "github.com/ose-micro/postman/internal/api/grpc/gen/go/ose/micro/postman/email/v1"
	"github.com/ose-micro/postman/internal/app"
	"github.com/ose-micro/postman/internal/business/email"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	EmailHandler struct {
		emailv1.UnimplementedEmailServiceServer
		app    email.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (e *EmailHandler) response(param email.Public) *emailv1.Email {
	return &emailv1.Email{
		Id:        param.Id,
		Count:     param.Count,
		Recipient: param.Recipient,
		Data:      stringifyInterfaceMap(param.Data),
		Subject:   param.Subject,
		Sender:    param.Sender,
		From:      param.From,
		Template:  param.Template,
		Message:   param.Message,
		State: func() emailv1.State {
			switch param.State {
			case email.StateComplete:
				return emailv1.State_StateComplete
			case email.StateFailed:
				return emailv1.State_StateFailed
			default:
				return emailv1.State_StateUnknown
			}
		}(),
		CreatedAt: timestamppb.New(param.CreatedAt),
		UpdatedAt: timestamppb.New(param.UpdatedAt),
	}
}

func (e *EmailHandler) Create(ctx context.Context, request *emailv1.CreateRequest) (*emailv1.CreateResponse, error) {
	ctx, span := e.tracer.Start(ctx, "api.grpc.email.create.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := email.CreateCommand{
		Recipient: request.Recipient,
		Sender:    request.Sender,
		Data:      convertStringMapToInterfaceMap(request.Data),
		Template:  request.Template,
		From:      request.From,
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

		return nil, parseError(err)
	}

	e.log.Info("email create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.Any("payload", request),
	)

	return &emailv1.CreateResponse{
		Message: "email create successfully",
		Record:  e.response(record.Public()),
	}, nil
}

func (e *EmailHandler) Resend(ctx context.Context, request *emailv1.ResendRequest) (*emailv1.ResendResponse, error) {
	ctx, span := e.tracer.Start(ctx, "api.grpc.email.resend.handler", trace.WithAttributes(
		attribute.String("operation", "resend"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := email.IdCommand{Id: request.Id}

	if err := e.app.Resend(ctx, payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to resend verification mail",
			zap.String("trace_id", traceId),
			zap.String("operation", "resend"),
			zap.Error(err),
		)

		return nil, parseError(err)
	}

	e.log.Info("resend verification mail process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "resend"),
		zap.Any("payload", request),
	)

	return &emailv1.ResendResponse{
		Message: "resend mail successfully",
	}, nil
}

func (e *EmailHandler) Read(ctx context.Context, request *emailv1.ReadRequest) (*emailv1.ReadResponse, error) {
	ctx, span := e.tracer.Start(ctx, "api.grpc.email.repository.handler", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	query, err := buildAppRequest(request.Request)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to case to dto",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)
		return nil, err
	}

	records, err := e.app.Read(ctx, *query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		e.log.Error("failed to repository roles",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)
		return nil, err
	}

	result := map[string]*emailv1.Emails{}

	for k, v := range records {
		switch x := v.(type) {
		case []email.Public:
			list := make([]*emailv1.Email, 0)
			for _, v := range x {
				list = append(list, e.response(v))
			}
			result[k] = &emailv1.Emails{
				Data: list,
			}
		}
	}

	return &emailv1.ReadResponse{
		Result: result,
	}, nil
}

func NewEmail(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *EmailHandler {
	return &EmailHandler{
		app:    apps.Email,
		log:    log,
		tracer: tracer,
	}
}

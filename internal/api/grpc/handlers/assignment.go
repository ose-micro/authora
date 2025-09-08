package handlers

import (
	"context"
	"fmt"

	assignmentv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/assignment/v1"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	AssignmentHandler struct {
		assignmentv1.UnimplementedAssignmentServiceServer
		app    assignment.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (h *AssignmentHandler) response(param assignment.Public) *assignmentv1.Assignment {
	return &assignmentv1.Assignment{
		Id:        param.Id,
		Tenant:    param.Tenant,
		User:      param.User,
		Role:      param.Role,
		Version:   param.Version,
		CreatedAt: timestamppb.New(param.CreatedAt),
		UpdatedAt: timestamppb.New(param.UpdatedAt),
		DeletedAt: buildDeletedAt(param.DeletedAt),
	}
}

func (h *AssignmentHandler) Create(ctx context.Context, request *assignmentv1.CreateRequest) (*assignmentv1.CreateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.assignment.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := assignment.CreateCommand{
		Role:   request.Role,
		Tenant: request.Tenant,
		User:   request.User,
	}

	record, err := h.app.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to create assignment",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("assignment create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("dto", request),
	)

	return &assignmentv1.CreateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *AssignmentHandler) Update(ctx context.Context, request *assignmentv1.UpdateRequest) (*assignmentv1.UpdateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.assignment.update.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := assignment.UpdateCommand{
		Id:   request.Id,
		Role: request.Role,
	}

	record, err := h.app.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update assignment",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("assignment update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	return &assignmentv1.UpdateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *AssignmentHandler) Read(ctx context.Context, request *assignmentv1.ReadRequest) (*assignmentv1.ReadResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.assignment.read.handler", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	query, err := buildAppRequest(request.Request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to case to dto",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)
		return nil, err
	}

	records, err := h.app.Read(ctx, assignment.ReadQuery{
		Request: *query,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to read organizations",
			zap.String("trace_id", traceId),
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	return &assignmentv1.ReadResponse{
		Result: func() map[string]*assignmentv1.Assignments {
			data := map[string]*assignmentv1.Assignments{}

			for k, v := range records {
				switch x := v.(type) {
				case []assignment.Public:
					list := make([]*assignmentv1.Assignment, 0)
					for _, v := range x {
						list = append(list, h.response(v))
					}
					data[k] = &assignmentv1.Assignments{
						Data: list,
					}
				}
			}

			return data
		}(),
	}, nil
}

func NewAssignment(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *AssignmentHandler {
	return &AssignmentHandler{
		app:    apps.Assignment,
		log:    log,
		tracer: tracer,
	}
}

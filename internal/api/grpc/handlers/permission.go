package handlers

import (
	"context"
	"fmt"

	permissionv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/permission/v1"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/permission"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	PermissionHandler struct {
		permissionv1.UnimplementedPermissionServiceServer
		app    permission.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (h *PermissionHandler) response(param permission.Public) *permissionv1.Permission {
	return &permissionv1.Permission{
		Id:        param.Id,
		Resource:  param.Resource,
		Action:    param.Action,
		Version:   param.Version,
		CreatedAt: timestamppb.New(param.CreatedAt),
		UpdatedAt: timestamppb.New(param.UpdatedAt),
		DeletedAt: buildDeletedAt(param.DeletedAt),
	}
}

func (h *PermissionHandler) Create(ctx context.Context, request *permissionv1.CreateRequest) (*permissionv1.CreateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.permission.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := permission.CreateCommand{
		Resource: request.Resource,
		Action:   request.Action,
	}

	record, err := h.app.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to create permission",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("permission create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("dto", request),
	)

	return &permissionv1.CreateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *PermissionHandler) Update(ctx context.Context, request *permissionv1.UpdateRequest) (*permissionv1.UpdateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.permission.update.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := permission.UpdateCommand{
		Id:       request.Id,
		Resource: request.Resource,
		Action:   request.Action,
	}

	record, err := h.app.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update permission",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("permission update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	return &permissionv1.UpdateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *PermissionHandler) Read(ctx context.Context, request *permissionv1.ReadRequest) (*permissionv1.ReadResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.permission.read.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
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
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	records, err := h.app.Read(ctx, permission.ReadQuery{
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

	return &permissionv1.ReadResponse{
		Result: func() map[string]*permissionv1.Permissions {
			data := map[string]*permissionv1.Permissions{}

			for k, v := range records {
				switch x := v.(type) {
				case []permission.Public:
					list := make([]*permissionv1.Permission, 0)
					for _, v := range x {
						list = append(list, h.response(v))
					}
					data[k] = &permissionv1.Permissions{
						Data: list,
					}
				}
			}

			return data
		}(),
	}, nil
}

func NewPermission(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *PermissionHandler {
	return &PermissionHandler{
		app:    apps.Permission,
		log:    log,
		tracer: tracer,
	}
}

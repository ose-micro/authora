package handlers

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/domain/role"
	rolev1 "github.com/ose-micro/authora/internal/interface/grpc/gen/go/ose/micro/authora/role/v1"
	commonv1 "github.com/ose-micro/authora/internal/interface/grpc/gen/go/ose/micro/common/v1"
	"github.com/ose-micro/common"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	RoleHandler struct {
		rolev1.UnimplementedRoleServiceServer
		app    role.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (h *RoleHandler) response(param role.Public) *rolev1.Role {
	return &rolev1.Role{
		Id:     param.Id,
		Name:   param.Name,
		Tenant: param.Tenant,
		Permissions: func() []*commonv1.Permission {
			var permissions []*commonv1.Permission

			for _, p := range param.Permissions {
				permissions = append(permissions, &commonv1.Permission{
					Resource: p.Resource,
					Action:   p.Action,
				})
			}

			return permissions
		}(),
		Version:   param.Version,
		CreatedAt: timestamppb.New(param.CreatedAt),
		UpdatedAt: timestamppb.New(param.UpdatedAt),
		DeletedAt: func() *timestamppb.Timestamp {
			if param.DeletedAt != nil {
				return timestamppb.New(*param.DeletedAt)
			}

			return nil
		}(),
	}
}

func (h *RoleHandler) Create(ctx context.Context, request *rolev1.CreateRequest) (*rolev1.CreateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "interface.grpc.role.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := role.CreateCommand{
		Name:        request.Name,
		Tenant:      request.Tenant,
		Description: request.Description,
		Permissions: func() []common.Permission {
			var permissions []common.Permission

			for _, p := range request.Permissions {
				permissions = append(permissions, common.Permission{
					Resource: p.Resource,
					Action:   p.Action,
				})
			}

			return permissions
		}(),
	}

	record, err := h.app.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to create role",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("role create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("payload", request),
	)

	return &rolev1.CreateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *RoleHandler) Update(ctx context.Context, request *rolev1.UpdateRequest) (*rolev1.UpdateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "interface.grpc.role.update.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := role.UpdateCommand{
		Id:          request.Id,
		Name:        request.Name,
		Tenant:      request.Tenant,
		Description: request.Description,
		Permissions: func() []common.Permission {
			var permissions []common.Permission

			for _, p := range request.Permissions {
				permissions = append(permissions, common.Permission{
					Resource: p.Resource,
					Action:   p.Action,
				})
			}

			return permissions
		}(),
	}

	record, err := h.app.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update role",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("role update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("payload", request),
	)

	return &rolev1.UpdateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *RoleHandler) Read(ctx context.Context, request *rolev1.ReadRequest) (*rolev1.ReadResponse, error) {
	ctx, span := h.tracer.Start(ctx, "interface.grpc.role.read.handler", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
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

	records, err := h.app.Read(ctx, role.ReadQuery{
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

	return &rolev1.ReadResponse{
		Result: func() map[string]*rolev1.Tenants {
			data := map[string]*rolev1.Tenants{}

			for k, v := range records {
				switch x := v.(type) {
				case []role.Public:
					list := make([]*rolev1.Role, 0)
					for _, v := range x {
						list = append(list, h.response(v))
					}
					data[k] = &rolev1.Tenants{
						Data: list,
					}
				}
			}

			return data
		}(),
	}, nil
}

func NewRole(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *RoleHandler {
	return &RoleHandler{
		app:    apps.Role,
		log:    log,
		tracer: tracer,
	}
}

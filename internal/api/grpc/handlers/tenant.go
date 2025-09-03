package handlers

import (
	"context"
	"fmt"

	tenantv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/tenant/v1"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/tenant"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	TenantHandler struct {
		tenantv1.UnimplementedTenantServiceServer
		app    tenant.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (h *TenantHandler) response(param tenant.Public) *tenantv1.Tenant {
	return &tenantv1.Tenant{
		Id:        param.Id,
		Name:      param.Name,
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

func (h *TenantHandler) Create(ctx context.Context, request *tenantv1.CreateRequest) (*tenantv1.CreateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.tenant.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := tenant.CreateCommand{
		Name: request.Name,
		Metadata: func() map[string]interface{} {
			if request.Metadata != nil {
				metadata := map[string]interface{}{}
				for k, v := range request.Metadata {
					metadata[k] = v
				}
			}

			return nil
		}(),
	}

	record, err := h.app.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to create tenant",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("tenant create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("payload", request),
	)

	return &tenantv1.CreateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *TenantHandler) Update(ctx context.Context, request *tenantv1.UpdateRequest) (*tenantv1.UpdateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.tenant.update.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("payload", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := tenant.UpdateCommand{
		Id:   request.Id,
		Name: request.Name,
	}

	record, err := h.app.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update tenant",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("tenant update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("payload", request),
	)

	return &tenantv1.UpdateResponse{
		Record: h.response(*record.Public()),
	}, nil
}

func (h *TenantHandler) Read(ctx context.Context, request *tenantv1.ReadRequest) (*tenantv1.ReadResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.tenant.read.handler", trace.WithAttributes(
		attribute.String("operation", "READ"),
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
			zap.String("operation", "READ"),
			zap.Error(err),
		)
		return nil, err
	}

	records, err := h.app.Read(ctx, tenant.ReadQuery{
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

	return &tenantv1.ReadResponse{
		Result: func() map[string]*tenantv1.Tenants {
			data := map[string]*tenantv1.Tenants{}

			for k, v := range records {
				switch x := v.(type) {
				case []tenant.Public:
					list := make([]*tenantv1.Tenant, 0)
					for _, v := range x {
						list = append(list, h.response(v))
					}
					data[k] = &tenantv1.Tenants{
						Data: list,
					}
				}
			}

			return data
		}(),
	}, nil
}

func NewTenant(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *TenantHandler {
	return &TenantHandler{
		app:    apps.Tenant,
		log:    log,
		tracer: tracer,
	}
}

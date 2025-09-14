package handlers

import (
	"context"
	"fmt"

	authv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/auth/v1"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/common/claims"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_error "github.com/ose-micro/error"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	AuthHandler struct {
		authv1.UnimplementedAuthServiceServer
		log    logger.Logger
		tracer tracing.Tracer
		app    user.App
	}
)

func (h *AuthHandler) HasRole(ctx context.Context, request *authv1.HasRoleRequest) (*authv1.HasRoleResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.auth.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.HasRoleCommand{
		Token:  request.Token,
		Role:   request.Role,
		Tenant: request.Tenant,
	}

	result, err := h.app.HasRole(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to has_role auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_role"),
			zap.Error(err),
		)

		return nil, parseError(err)
	}

	if !result {
		err := ose_error.New(ose_error.ErrUnauthorized, "does not have role")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to has_role auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_role"),
			zap.Error(err),
		)

		return nil, parseError(err)
	}

	return &authv1.HasRoleResponse{
		Message: "has role",
	}, nil
}

func (h *AuthHandler) HasPermission(ctx context.Context, request *authv1.HasPermissionRequest) (*authv1.HasPermissionResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.auth.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.HasPermissionCommand{
		Token:  request.Token,
		Tenant: request.Tenant,
		Permission: func() *claims.Permission {
			if request.Permission == nil {
				return nil
			}

			return &claims.Permission{
				Action:   request.Permission.Action,
				Resource: request.Permission.Resource,
			}
		}(),
	}

	result, err := h.app.HasPermission(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to has_permission auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_permission"),
			zap.Error(err),
		)

		return nil, err
	}

	if !result {
		err := ose_error.New(ose_error.ErrUnauthorized, "does not have permission")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to has_permission auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_permission"),
			zap.Error(err),
		)

		return nil, err
	}

	return &authv1.HasPermissionResponse{
		Message: "has permission",
	}, nil
}

func (h *AuthHandler) RequestPurposeToken(ctx context.Context, request *authv1.RequestPurposeTokenRequest) (*authv1.RequestPurposeTokenResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.auth.request_purpose_token.handler", trace.WithAttributes(
		attribute.String("operation", "request_purpose_token"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.PurposeTokenCommand{
		Id:      request.Id,
		Purpose: request.Purpose,
	}

	token, err := h.app.RequestPurposeToken(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to request_purpose_token auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_purpose_token"),
			zap.Error(err),
		)

		return nil, err
	}

	return &authv1.RequestPurposeTokenResponse{
		Token: *token,
	}, nil
}

func (h *AuthHandler) RequestAccessToken(ctx context.Context, request *authv1.RequestAccessTokenRequest) (*authv1.RequestAccessTokenResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.auth.request_access_token.handler", trace.WithAttributes(
		attribute.String("operation", "request_access_token"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.TokenCommand{
		Token: request.Refresh,
	}

	token, err := h.app.RequestAccessToken(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to request_access_token auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	return &authv1.RequestAccessTokenResponse{
		Token: *token,
	}, nil
}

func (h *AuthHandler) ParseClaim(ctx context.Context, request *authv1.ParseClaimRequest) (*authv1.ParseClaimResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.auth.parse_claim.handler", trace.WithAttributes(
		attribute.String("operation", "parse_claim"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.TokenCommand{
		Token: request.Token,
	}

	claim, err := h.app.ParseClaims(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to parse_claim auth",
			zap.String("trace_id", traceId),
			zap.String("operation", "parse_claim"),
			zap.Error(err),
		)

		return nil, err
	}

	return &authv1.ParseClaimResponse{
		Message: "parse_claim auth",
		Claim:   buildClaim(*claim),
	}, nil
}

func NewAuth(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *AuthHandler {
	return &AuthHandler{
		app:    apps.User,
		log:    log,
		tracer: tracer,
	}
}

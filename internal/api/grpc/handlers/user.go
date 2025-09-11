package handlers

import (
	"context"
	"fmt"

	userv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/user/v1"
	commonv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/common/v1"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/common"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_error "github.com/ose-micro/error"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	UserHandler struct {
		userv1.UnimplementedUserServiceServer
		app    user.App
		log    logger.Logger
		tracer tracing.Tracer
	}
)

func (h *UserHandler) response(param user.Public) (*userv1.User, error) {
	metadata := map[string]string{}

	if param.Metadata != nil {
		meta, err := common.ToStringMap(param.Metadata)
		if err != nil {
			return nil, ose_error.New(ose_error.ErrInternal, err.Error())
		}

		metadata = meta
	}

	return &userv1.User{
		Id:         param.Id,
		GivenNames: param.GivenNames,
		FamilyName: param.FamilyName,
		Email:      param.Email,
		Password:   param.Password,
		Metadata:   metadata,
		Version:    param.Version,
		Status:     buildUserStatus(*param.Status),
		CreatedAt:  timestamppb.New(param.CreatedAt),
		UpdatedAt:  timestamppb.New(param.UpdatedAt),
		DeletedAt:  buildDeletedAt(param.DeletedAt),
	}, nil
}

func (h *UserHandler) Create(ctx context.Context, request *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.create.handler", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.CreateCommand{
		GivenNames: request.GivenNames,
		FamilyName: request.FamilyName,
		Email:      request.Email,
		Password:   request.Password,
		Role:       request.Role,
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
		h.log.Error("failed to create user",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("user create process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "create"),
		zap.Any("dto", request),
	)

	result, err := h.response(*record.Public())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to create user",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err))
		return nil, err
	}
	return &userv1.CreateResponse{
		Record: result,
	}, nil
}

func (h *UserHandler) Update(ctx context.Context, request *userv1.UpdateRequest) (*userv1.UpdateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.update.handler", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.UpdateCommand{
		Id:         request.Id,
		GivenNames: request.GivenNames,
		FamilyName: request.FamilyName,
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

	record, err := h.app.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update user",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("user update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	result, err := h.response(*record.Public())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update user",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err))

		return nil, err
	}

	return &userv1.UpdateResponse{
		Record: result,
	}, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, request *userv1.ChangePasswordRequest) (*userv1.ChangePasswordResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.change_password.handler", trace.WithAttributes(
		attribute.String("operation", "change_password"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.ChangePasswordCommand{
		Id:          request.Id,
		OldPassword: request.OldPassword,
		Password:    request.Password,
	}

	record, err := h.app.ChangePassword(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to update user",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("user update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	result, _ := h.response(*record.Public())

	return &userv1.ChangePasswordResponse{
		Record: result,
	}, nil
}

func (h *UserHandler) RequestAccessToken(ctx context.Context, request *userv1.RequestAccessTokenRequest) (*userv1.RequestAccessTokenResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.request_access_token.handler", trace.WithAttributes(
		attribute.String("operation", "login"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.TokenCommand{
		Token: request.Token,
	}

	res, err := h.app.RequestAccessToken(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to request access token",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("request access token process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	return &userv1.RequestAccessTokenResponse{
		Token: *res,
	}, nil
}

func (h *UserHandler) RequestPurposeToken(ctx context.Context, request *userv1.RequestPurposeTokenRequest) (*userv1.RequestPurposeTokenResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.request_purpose_token.handler", trace.WithAttributes(
		attribute.String("operation", "request_purpose_token"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.PurposeTokenCommand{
		Purpose: request.Purpose,
		Id:      request.Id,
	}

	res, err := h.app.RequestPurposeToken(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to request_purpose_token user",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_purpose_token"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("user update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	return &userv1.RequestPurposeTokenResponse{
		Token: *res,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.login.handler", trace.WithAttributes(
		attribute.String("operation", "login"),
		attribute.String("dto", fmt.Sprintf("%v", request)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	payload := user.LoginCommand{
		Email:    request.Email,
		Password: request.Password,
	}

	res, err := h.app.Login(ctx, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("failed to login user",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	h.log.Info("user update process successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "update"),
		zap.Any("dto", request),
	)

	return &userv1.LoginResponse{
		Message: "success",
		Record: &commonv1.Auth{
			Access:  res.Access,
			Refresh: res.Refresh,
		},
	}, nil
}

func (h *UserHandler) Read(ctx context.Context, request *userv1.ReadRequest) (*userv1.ReadResponse, error) {
	ctx, span := h.tracer.Start(ctx, "api.grpc.user.read.handler", trace.WithAttributes(
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

	records, err := h.app.Read(ctx, user.ReadQuery{
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

	return &userv1.ReadResponse{
		Result: func() map[string]*userv1.Users {
			data := map[string]*userv1.Users{}

			for k, v := range records {
				switch x := v.(type) {
				case []user.Public:
					list := make([]*userv1.User, 0)
					for _, v := range x {
						result, err := h.response(v)
						if err != nil {

						}
						list = append(list, result)
					}
					data[k] = &userv1.Users{
						Data: list,
					}
				}
			}

			return data
		}(),
	}, nil
}

func NewUser(apps app.Apps, log logger.Logger, tracer tracing.Tracer) *UserHandler {
	return &UserHandler{
		app:    apps.User,
		log:    log,
		tracer: tracer,
	}
}

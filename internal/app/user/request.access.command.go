package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
	"github.com/ose-micro/common"
	"github.com/ose-micro/common/claims"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	ose_jwt "github.com/ose-micro/jwt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type requestAccessTokenCommandHandler struct {
	log    logger.Logger
	repo   repository.Repository
	tracer tracing.Tracer
	jwt    ose_jwt.Manager
	cache  user.Cache
}

// Handle implements cqrs.CommandHandle.
func (h requestAccessTokenCommandHandler) Handle(ctx context.Context, command user.TokenCommand) (*string, error) {
	ctx, span := h.tracer.Start(ctx, "app.user.request_access_token.command.handler", trace.WithAttributes(
		attribute.String("operation", "request_access_token"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command dto
	if err := command.Validate(); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrBadRequest, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	tokenClaim, err := h.cache.Get(ctx, command.Token)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrUnauthorized, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	token, err := h.prepareToken(ctx, *tokenClaim)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrUnauthorized, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	return token, nil
}

func (h requestAccessTokenCommandHandler) prepareToken(ctx context.Context, payload user.Token) (*string, error) {
	tenants := make(map[string]ose_jwt.Tenant)

	// Fetch all roles for user
	assignsFacts, err := h.repo.Assignment.Read(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "user",
						Op:    dto.OpEq,
						Value: payload.User(),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	fact := assignsFacts["one"]
	var assigns []assignment.Public

	if err := common.JsonToAny(fact, &assigns); err != nil {
		return nil, err
	}

	for _, assign := range assigns {
		one, err := h.repo.Role.ReadOne(ctx, dto.Request{
			Queries: []dto.Query{
				{
					Name: "one",
					Filters: []dto.Filter{
						{
							Field: "_id",
							Op:    dto.OpEq,
							Value: assign.Role,
						},
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		permissions, err := h.preparePermission(ctx, *one)
		if err != nil {
			return nil, err
		}

		tenants[assign.Tenant] = ose_jwt.Tenant{
			Role:        one.ID(),
			Tenant:      one.Tenant(),
			Permissions: permissions,
		}
	}

	token, _, err := h.jwt.IssueAccessToken(payload.User(), tenants, nil)
	if err != nil {
		return nil, err
	}
	accessToken, err := user.NewToken(user.TokenParam{
		User:    payload.User(),
		Purpose: "access",
		Tenant:  payload.Tenant(),
		Token:   token,
	})
	if err != nil {
		return nil, err
	}

	token = accessToken.Key()

	if err := h.cache.Save(ctx, accessToken, 15*time.Minute); err != nil {
		return nil, err
	}

	return &token, nil
}

func (h requestAccessTokenCommandHandler) preparePermission(ctx context.Context, one role.Domain) ([]claims.Permission, error) {

	list := make([]claims.Permission, 0)

	for _, id := range one.Permissions() {
		permission, err := h.repo.Permission.ReadOne(ctx, dto.Request{
			Queries: []dto.Query{
				{
					Name: "one",
					Filters: []dto.Filter{
						{
							Field: "_id",
							Op:    dto.OpEq,
							Value: id,
						},
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		list = append(list, claims.Permission{
			Resource: permission.Resource(),
			Action:   permission.Action(),
		})
	}

	return list, nil
}

func newRequestAccessTokenCommandHandler(log logger.Logger, tracer tracing.Tracer,
	jwt ose_jwt.Manager, repo repository.Repository, cache user.Cache) cqrs.CommandHandle[user.TokenCommand, *string] {
	return &requestAccessTokenCommandHandler{
		log:    log,
		tracer: tracer,
		jwt:    jwt,
		repo:   repo,
		cache:  cache,
	}
}

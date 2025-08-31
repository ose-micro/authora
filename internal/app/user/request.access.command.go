package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/domain/assignment"
	"github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/common"
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
}

// Handle implements cqrs.CommandHandle.
func (h requestAccessTokenCommandHandler) Handle(ctx context.Context, command user.TokenCommand) (*string, error) {
	ctx, span := h.tracer.Start(ctx, "app.user.request_access_token.command.handler", trace.WithAttributes(
		attribute.String("operation", "request_access_token"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command payload
	if err := command.Validate(); err != nil {
		err := ose_error.New(ose_error.ErrInvalidInput, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	claims, err := h.jwt.ParseClaims(command.Token)
	if err != nil {
		err := ose_error.New(ose_error.ErrUnauthorized, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_access_token"),
			zap.Error(err),
		)

		return nil, err
	}

	token, err := h.prepareToken(ctx, claims.Sub)
	if err != nil {
		err := ose_error.New(ose_error.ErrUnauthorized, err.Error())
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

func (h requestAccessTokenCommandHandler) prepareToken(ctx context.Context, id string) (*string, error) {
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
						Value: id,
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

		tenants[assign.Tenant] = ose_jwt.Tenant{
			Role:        one.ID(),
			Tenant:      one.Tenant(),
			Permissions: one.Permissions(),
		}
	}

	token, _, err := h.jwt.IssueAccessToken(id, tenants, nil)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func newRequestAccessTokenCommandHandler(log logger.Logger, tracer tracing.Tracer,
	jwt ose_jwt.Manager, repo repository.Repository) cqrs.CommandHandle[user.TokenCommand, *string] {
	return &requestAccessTokenCommandHandler{
		log:    log,
		tracer: tracer,
		jwt:    jwt,
		repo:   repo,
	}
}

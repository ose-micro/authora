package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/repository"
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

type requestPurposeTokenCommandHandler struct {
	log    logger.Logger
	repo   repository.Repository
	tracer tracing.Tracer
	jwt    ose_jwt.Manager
}

// Handle implements cqrs.CommandHandle.
func (h requestPurposeTokenCommandHandler) Handle(ctx context.Context, command user.PurposeTokenCommand) (*string, error) {
	ctx, span := h.tracer.Start(ctx, "app.user.request_purpose_token.command.handler", trace.WithAttributes(
		attribute.String("operation", "request_purpose_token"),
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
			zap.String("operation", "request_purpose_token"),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := h.repo.User.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "_id",
						Op:    dto.OpEq,
						Value: command.Id,
					},
				},
			},
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_purpose_token"),
			zap.Error(err),
		)

		return nil, err
	}

	if !command.Safe {
		if !record.Status().IsActive() {
			err = ose_error.New(ose_error.ErrUnauthorized, "user is not active")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			h.log.Error("validation process failed",
				zap.String("trace_id", traceId),
				zap.String("operation", "request_purpose_token"),
				zap.Error(err),
			)

			return nil, err
		}
	}

	token, err := h.prepareToken(ctx, command.Id, command.Purpose)
	if err != nil {
		err := ose_error.New(ose_error.ErrUnauthorized, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "request_purpose_token"),
			zap.Error(err),
		)

		return nil, err
	}

	return token, nil
}

func (h requestPurposeTokenCommandHandler) prepareToken(ctx context.Context, id, purpose string) (*string, error) {
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

	sub := fmt.Sprintf("%s:%s", id, purpose)
	token, _, err := h.jwt.IssuePurposeToken(sub, tenants, nil)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (h requestPurposeTokenCommandHandler) preparePermission(ctx context.Context, one role.Domain) ([]claims.Permission, error) {

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

func newRequestPurposeTokenCommandHandler(log logger.Logger, tracer tracing.Tracer,
	jwt ose_jwt.Manager, repo repository.Repository) cqrs.CommandHandle[user.PurposeTokenCommand, *string] {
	return &requestPurposeTokenCommandHandler{
		log:    log,
		tracer: tracer,
		jwt:    jwt,
		repo:   repo,
	}
}

package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/common"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/core/utils"
	"github.com/ose-micro/cqrs"
	ose_error "github.com/ose-micro/error"
	ose_jwt "github.com/ose-micro/jwt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginCommandHandler struct {
	repo   repository.Repository
	log    logger.Logger
	tracer tracing.Tracer
	bs     business.Domain
	jwt    ose_jwt.Manager
}

// Handle implements cqrs.CommandHandle.
func (u loginCommandHandler) Handle(ctx context.Context, command user.LoginCommand) (*user.Auth, error) {
	ctx, span := u.tracer.Start(ctx, "app.user.login.command.handler", trace.WithAttributes(
		attribute.String("operation", "login"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command dto
	if err := command.Validate(); err != nil {
		err := ose_error.New(ose_error.ErrInvalidInput, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	record, err := u.repo.User.ReadOne(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "email",
						Op:    dto.OpEq,
						Value: command.Email,
					},
				},
			},
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	if !record.Status().IsActive() {
		msg := fmt.Sprintf("user is %s, user need to be activated", strings.ToLower(record.Status().State.String()))
		err := ose_error.New(ose_error.ErrUnauthorized, msg)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	if !utils.CheckPasswordHash(command.Password, record.Password()) {
		err := ose_error.New(ose_error.ErrUnauthorized, "password does not match")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to login record",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	auth, err := u.prepareAuth(ctx, *record)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.log.Error("failed to login record",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)

		return nil, err
	}

	u.log.Info("login process complete successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "login"),
		zap.Any("dto", command),
	)

	return auth, nil
}

func (u loginCommandHandler) prepareAuth(ctx context.Context, command user.Domain) (*user.Auth, error) {
	tenants := make(map[string]ose_jwt.Tenant)

	// Fetch all roles for user
	assignsFacts, err := u.repo.Assignment.Read(ctx, dto.Request{
		Queries: []dto.Query{
			{
				Name: "one",
				Filters: []dto.Filter{
					{
						Field: "user",
						Op:    dto.OpEq,
						Value: command.ID(),
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
		one, err := u.repo.Role.ReadOne(ctx, dto.Request{
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

		permissions, err := u.preparePermission(ctx, *one)
		if err != nil {
			return nil, err
		}

		tenants[assign.Tenant] = ose_jwt.Tenant{
			Role:        one.ID(),
			Tenant:      one.Tenant(),
			Permissions: permissions,
		}
	}

	accessToken, _, err := u.jwt.IssueAccessToken(command.ID(), tenants, nil)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := u.jwt.IssueRefreshToken(command.ID(), tenants, nil)

	return &user.Auth{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (u loginCommandHandler) preparePermission(ctx context.Context, one role.Domain) ([]common.Permission, error) {

	list := make([]common.Permission, 0)

	for _, id := range one.Permissions() {
		permission, err := u.repo.Permission.ReadOne(ctx, dto.Request{
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

		list = append(list, common.Permission{
			Resource: permission.Resource(),
			Action:   permission.Action(),
		})
	}

	return list, nil
}

func newLoginCommandHandler(bs business.Domain, repo repository.Repository,
	log logger.Logger, tracer tracing.Tracer, jwt ose_jwt.Manager) cqrs.CommandHandle[user.LoginCommand, *user.Auth] {
	return &loginCommandHandler{
		repo:   repo,
		log:    log,
		tracer: tracer,
		bs:     bs,
		jwt:    jwt,
	}
}

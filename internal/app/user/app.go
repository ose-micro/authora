package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/infrastruture/cache"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	ose_jwt "github.com/ose-micro/jwt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type app struct {
	tracer              tracing.Tracer
	log                 logger.Logger
	create              cqrs.CommandHandle[user.CreateCommand, *user.Domain]
	update              cqrs.CommandHandle[user.UpdateCommand, *user.Domain]
	login               cqrs.CommandHandle[user.LoginCommand, *user.Auth]
	logout              cqrs.CommandHandle[user.TokenCommand, bool]
	hasRole             cqrs.CommandHandle[user.HasRoleCommand, *bool]
	changeStatus        cqrs.CommandHandle[user.StatusCommand, *bool]
	hasPermission       cqrs.CommandHandle[user.HasPermissionCommand, *bool]
	parseClaims         cqrs.CommandHandle[user.TokenCommand, *user.TokenClaim]
	requestPurposeToken cqrs.CommandHandle[user.PurposeTokenCommand, *string]
	requestAccessToken  cqrs.CommandHandle[user.TokenCommand, *string]
	changePassword      cqrs.CommandHandle[user.ChangePasswordCommand, *user.Domain]
	resetPassword       cqrs.CommandHandle[user.ResetPasswordCommand, *user.Domain]
	read                cqrs.QueryHandle[user.ReadQuery, map[string]any]
	readOne             cqrs.QueryHandle[user.ReadQuery, *user.Domain]
}

// Logout implements user.App.
func (a *app) Logout(ctx context.Context, command user.TokenCommand) (bool, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.logout.command", trace.WithAttributes(
		attribute.String("operation", "logout"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := a.logout.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "logout"),
			zap.Error(err),
		)
		return false, err
	}

	return record, nil
}

func (a app) ResetPassword(ctx context.Context, command user.ResetPasswordCommand) (*user.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.reset_password.query", trace.WithAttributes(
		attribute.String("operation", "reset_password"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := a.resetPassword.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "reset_password"),
			zap.Error(err),
		)
		return nil, err
	}

	return record, nil
}

func (a app) ReadOne(ctx context.Context, command user.ReadQuery) (*user.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.read_one.query", trace.WithAttributes(
		attribute.String("operation", "read_one"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	record, err := a.readOne.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)
		return nil, err
	}

	return record, nil
}

func (a app) ChangeStatus(ctx context.Context, command user.StatusCommand) (bool, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.change_status.command", trace.WithAttributes(
		attribute.String("operation", "change_status"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := a.changeStatus.Handle(ctx, command); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_status"),
			zap.Error(err),
		)
		return false, err
	}

	return true, nil
}

func (a app) RequestAccessToken(ctx context.Context, command user.TokenCommand) (*string, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.request.access.token.command", trace.WithAttributes(
		attribute.String("operation", "access_token"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.requestAccessToken.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "access_token"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) RequestPurposeToken(ctx context.Context, command user.PurposeTokenCommand) (*string, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.request.purpose.token.command", trace.WithAttributes(
		attribute.String("operation", "purpose_token"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.requestPurposeToken.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "purpose_token"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) HasRole(ctx context.Context, command user.HasRoleCommand) (bool, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.has_role.command", trace.WithAttributes(
		attribute.String("operation", "has_role"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.hasRole.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_role"),
			zap.Error(err),
		)
		return false, err
	}

	return *res, nil
}

func (a app) HasPermission(ctx context.Context, command user.HasPermissionCommand) (bool, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.has_permission.command", trace.WithAttributes(
		attribute.String("operation", "has_permission"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.hasPermission.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "has_permission"),
			zap.Error(err),
		)
		return false, err
	}

	return *res, nil
}

func (a app) ParseClaims(ctx context.Context, command user.TokenCommand) (*user.TokenClaim, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.parse_claims.command", trace.WithAttributes(
		attribute.String("operation", "parse_claims"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.parseClaims.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "parse_claims"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) Login(ctx context.Context, command user.LoginCommand) (*user.Auth, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.login.command", trace.WithAttributes(
		attribute.String("operation", "login"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.login.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "login"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) ChangePassword(ctx context.Context, command user.ChangePasswordCommand) (*user.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.change_password.command", trace.WithAttributes(
		attribute.String("operation", "change_password"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.changePassword.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "change_password"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) Read(ctx context.Context, command user.ReadQuery) (map[string]any, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.repository.command", trace.WithAttributes(
		attribute.String("operation", "repository"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.read.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "repository"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) Create(ctx context.Context, command user.CreateCommand) (*user.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.create.command", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.create.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) Update(ctx context.Context, command user.UpdateCommand) (*user.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.update.command", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.update.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "update"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) Delete(ctx context.Context, params user.UpdateCommand) (*user.Domain, error) {
	//TODO implement me
	panic("implement me")
}

func NewApp(bs business.Domain, log logger.Logger, tracer tracing.Tracer, repo repository.Repository,
	jwt ose_jwt.Manager, bus domain.Bus, cache cache.Cache) user.App {
	return &app{
		tracer:              tracer,
		log:                 log,
		create:              newCreateCommandHandler(bs, repo, log, tracer, bus),
		update:              newUpdateCommandHandler(bs, repo, log, tracer),
		login:               newLoginCommandHandler(bs, repo, log, tracer, jwt, cache.Token),
		hasRole:             newHasRoleCommandHandler(log, tracer, jwt),
		changeStatus:        newChangeStausCommandHandler(bs, repo, log, tracer),
		hasPermission:       newHasPermissionCommandHandler(log, tracer, jwt),
		parseClaims:         newParseClaimCommandHandler(log, tracer, jwt, cache.Token),
		requestPurposeToken: newRequestPurposeTokenCommandHandler(log, tracer, jwt, repo, cache.Token),
		requestAccessToken:  newRequestAccessTokenCommandHandler(log, tracer, jwt, repo, cache.Token),
		changePassword:      newChangePasswordCommandHandler(bs, repo, log, tracer),
		read:                newReadQueryHandler(repo.User, log, tracer),
		readOne:             newReadOneQueryHandler(repo.User, log, tracer),
		resetPassword:       newResetPasswordCommandHandler(bs, repo, log, tracer),
		logout:              newLogoutCommandHandler(log, tracer, cache.Token),
	}
}

package user

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/domain/user"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type app struct {
	tracer         tracing.Tracer
	log            logger.Logger
	create         cqrs.CommandHandle[user.CreateCommand, *user.Domain]
	update         cqrs.CommandHandle[user.UpdateCommand, *user.Domain]
	login          cqrs.CommandHandle[user.LoginCommand, *user.Auth]
	changePassword cqrs.CommandHandle[user.ChangePasswordCommand, *user.Domain]
	read           cqrs.QueryHandle[user.ReadQuery, map[string]any]
}

func (a app) Login(ctx context.Context, command user.LoginCommand) (*user.Auth, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.login.command", trace.WithAttributes(
		attribute.String("operation", "login"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
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
		attribute.String("payload", fmt.Sprintf("%v", command)),
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
	ctx, span := a.tracer.Start(ctx, "app.user.read.command", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.read.Handle(ctx, command)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		a.log.Error("failed to process command",
			zap.String("trace_id", traceId),
			zap.String("operation", "read"),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (a app) Create(ctx context.Context, command user.CreateCommand) (*user.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.user.create.command", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
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
		attribute.String("payload", fmt.Sprintf("%v", command)),
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

func NewApp(bs domain.Domain, log logger.Logger, tracer tracing.Tracer, repo repository.Repository) user.App {
	return &app{
		tracer:         tracer,
		log:            log,
		create:         newCreateCommandHandler(bs, repo, log, tracer),
		update:         newUpdateCommandHandler(bs, repo, log, tracer),
		read:           newReadQueryHandler(repo.User, log, tracer),
		changePassword: newChangePasswordCommandHandler(bs, repo, log, tracer),
		login:          newLoginCommandHandler(bs, repo, log, tracer),
	}
}

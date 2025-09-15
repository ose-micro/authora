package role

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/role"
	"github.com/ose-micro/authora/internal/repository"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type app struct {
	tracer tracing.Tracer
	log    logger.Logger
	create domain.CommandHandle[role.CreateCommand, *role.Domain]
	update domain.CommandHandle[role.UpdateCommand, *role.Domain]
	read   domain.QueryHandle[role.ReadQuery, map[string]any]
}

func (a app) Read(ctx context.Context, command role.ReadQuery) (map[string]any, error) {
	ctx, span := a.tracer.Start(ctx, "app.role.repository.command", trace.WithAttributes(
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

func (a app) Create(ctx context.Context, command role.CreateCommand) (*role.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.role.create.command", trace.WithAttributes(
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

func (a app) Update(ctx context.Context, command role.UpdateCommand) (*role.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.role.update.command", trace.WithAttributes(
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

func (a app) Delete(ctx context.Context, params role.UpdateCommand) (*role.Domain, error) {
	//TODO implement me
	panic("implement me")
}

func NewApp(bs business.Domain, log logger.Logger, tracer tracing.Tracer, repo repository.Repository) role.App {
	return &app{
		tracer: tracer,
		log:    log,
		create: newCreateCommandHandler(bs, repo, log, tracer),
		update: newUpdateCommandHandler(bs, repo, log, tracer),
		read:   newReadQueryHandler(repo.Role, log, tracer),
	}
}

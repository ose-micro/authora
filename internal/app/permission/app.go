package permission

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/permission"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type app struct {
	tracer  tracing.Tracer
	log     logger.Logger
	create  domain.CommandHandle[permission.CreateCommand, *permission.Domain]
	update  domain.CommandHandle[permission.UpdateCommand, *permission.Domain]
	read    domain.QueryHandle[permission.ReadQuery, map[string]any]
	readOne domain.QueryHandle[permission.ReadQuery, *permission.Domain]
}

func (a app) ReadOne(ctx context.Context, command permission.ReadQuery) (*permission.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.permission.read_one.command", trace.WithAttributes(
		attribute.String("operation", "read_one"),
		attribute.String("dto", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := a.readOne.Handle(ctx, command)
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

	return res, nil
}

func (a app) Read(ctx context.Context, command permission.ReadQuery) (map[string]any, error) {
	ctx, span := a.tracer.Start(ctx, "app.permission.repository.command", trace.WithAttributes(
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

func (a app) Create(ctx context.Context, command permission.CreateCommand) (*permission.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.permission.create.command", trace.WithAttributes(
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

func (a app) Update(ctx context.Context, command permission.UpdateCommand) (*permission.Domain, error) {
	ctx, span := a.tracer.Start(ctx, "app.permission.update.command", trace.WithAttributes(
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

func (a app) Delete(ctx context.Context, params permission.UpdateCommand) (*permission.Domain, error) {
	//TODO implement me
	panic("implement me")
}

func NewApp(bs business.Domain, log logger.Logger, tracer tracing.Tracer, repo repository.Repository) permission.App {
	return &app{
		tracer:  tracer,
		log:     log,
		create:  newCreateCommandHandler(bs, repo, log, tracer),
		update:  newUpdateCommandHandler(bs, repo, log, tracer),
		read:    newReadQueryHandler(repo.Permission, log, tracer),
		readOne: newReadOneQueryHandler(repo.Permission, log, tracer),
	}
}

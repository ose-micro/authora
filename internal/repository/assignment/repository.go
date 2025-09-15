package assignment

import (
	"context"
	"fmt"

	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/business/assignment"
	"github.com/ose-micro/common"
	"github.com/ose-micro/core/dto"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_error "github.com/ose-micro/error"
	mongodb "github.com/ose-micro/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repository struct {
	collection *mongo.Collection
	log        logger.Logger
	tracer     tracing.Tracer
	bs         business.Domain
}

func (r *repository) Delete(ctx context.Context, payload assignment.Domain) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) ReadOne(ctx context.Context, request dto.Request) (*assignment.Domain, error) {
	ctx, span := r.tracer.Start(ctx, "repository.infrastructure.assignment.read_one", trace.WithAttributes(
		attribute.String("operation", "read_one"),
		attribute.String("dto", fmt.Sprintf("%v", request))),
	)
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	res, err := r.Read(ctx, request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository res",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)
		return nil, err
	}

	raw, ok := res["one"]
	if !ok {
		return nil, ose_error.New(ose_error.ErrNotFound, "assignment not found")
	}

	var records []assignment.Public

	if err := common.JsonToAny(raw, &records); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to repository res",
			zap.String("trace_id", traceId),
			zap.String("operation", "read_one"),
			zap.Error(err),
		)
		return nil, err
	}

	if len(records) == 0 {
		return nil, ose_error.New(ose_error.ErrNotFound, "assignment not found", traceId)
	}

	return r.toDomain(records[0]), nil
}

// Create implements assignment.Repository.
func (r *repository) Create(ctx context.Context, payload assignment.Domain) error {
	ctx, span := r.tracer.Start(ctx, "repository.infrastructure.assignment.create", trace.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("dto", fmt.Sprintf("%v", payload.Public())),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	record := newCollection(payload)
	if _, err := r.collection.InsertOne(ctx, record); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("failed to create in mongo",
			zap.String("trace_id", traceId),
			zap.String("operation", "create"),
			zap.Error(err),
		)
		return err
	}

	r.log.Info("create process complete successfully",
		zap.String("operation", "create"),
		zap.String("trace_id", traceId),
		zap.Any("dto", payload.Public()),
	)
	return nil
}

// Read implements assignment.Repository.
func (r *repository) Read(ctx context.Context, request dto.Request) (map[string]any, error) {
	ctx, span := r.tracer.Start(ctx, "repository.infrastructure.assignment.repository", trace.WithAttributes(
		attribute.String("operation", "read"),
		attribute.String("dto", fmt.Sprintf("%+v", request)),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	mongodb.RegisterType("assignment", assignment.Public{})
	typeHints := map[string]string{}

	for _, v := range request.Queries {
		typeHints[v.Name] = "assignment"
	}

	res, err := mongodb.RunFaceted(ctx, r.collection, request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.log.Error("Failed to fetch assignment by request",
			zap.String("operation", "repository"),
			zap.String("trace_id", traceID),
			zap.Any("dto", request),
			zap.Error(err),
		)
		return nil, ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceID)
	}

	r.log.Info("Read process completed successfully",
		zap.String("operation", "repository"),
		zap.String("trace_id", traceID),
		zap.Any("dto", request),
	)

	records, err := mongodb.CastFacetedResult(res, typeHints)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceID)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.log.Error("Failed to cast faceted result",
			zap.String("operation", "repository"),
			zap.String("trace_id", traceID),
			zap.Any("dto", request),
			zap.Error(err),
		)
		return nil, err
	}

	return records, nil
}

// Update implements assignment.Repository.
func (r *repository) Update(ctx context.Context, payload assignment.Domain) error {
	ctx, span := r.tracer.Start(ctx, "infrastructure.repository.assignment.update", trace.WithAttributes(
		attribute.String("operation", "update"),
		attribute.String("dto", fmt.Sprintf("%+v", payload.Public())),
	))
	defer span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	collection := newCollection(payload)
	filter := bson.M{"_id": payload.ID()}

	if _, err := r.collection.UpdateOne(ctx, filter, bson.M{
		"$set": collection,
	}); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceID)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.log.Error("failed to update assignment",
			zap.String("operation", "update"),
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		return err
	}

	r.log.Info("update process complete successfully",
		zap.String("operation", "update"),
		zap.String("trace_id", traceID),
		zap.Any("dto", payload.Public()),
	)

	return nil
}

func (r *repository) toDomain(payload assignment.Public) *assignment.Domain {
	result, _ := r.bs.Assignment.Existing(*payload.Params())
	return result
}

func NewRepository(db *mongodb.Client, log logger.Logger, tracer tracing.Tracer, bs business.Domain) assignment.Repo {
	return &repository{
		log:        log,
		tracer:     tracer,
		bs:         bs,
		collection: db.Collection("assignments"),
	}
}

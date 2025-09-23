package token

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	ose_error "github.com/ose-micro/error"
	"github.com/ose-micro/redis"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenCache struct {
	redis *redis.Client
	log   logger.Logger
	trace tracing.Tracer
}

func (t tokenCache) Save(ctx context.Context, payload *user.Token, ttl time.Duration) error {
	ctx, span := t.trace.Start(ctx, "user.cache.save", trace.WithAttributes(
		attribute.String("operation", "save"),
		attribute.String("dto", fmt.Sprintf("%v", payload))),
	)
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	data, err := json.Marshal(payload.Param())
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("error marshalling payload",
			zap.String("operation", "save"),
			zap.String("trace", traceId),
			zap.Error(err),
		)
		return err
	}

	err = t.redis.Set(ctx, payload.Key(), data, ttl)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("redis set error",
			zap.String("operation", "save"),
			zap.String("trace", traceId),
			zap.Error(err),
		)
		return err
	}

	t.log.Info("create process complete successfully",
		zap.String("operation", "create"),
		zap.String("trace_id", traceId),
		zap.Any("dto", fmt.Sprintf("%v", payload.Param())),
	)
	return nil
}

func (t tokenCache) Get(ctx context.Context, key string) (*user.Token, error) {
	ctx, span := t.trace.Start(ctx, "user.cache.get", trace.WithAttributes(
		attribute.String("operation", "get"),
		attribute.String("dto", key)),
	)
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := t.redis.Get(ctx, key)
	if err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("redis set error",
			zap.String("operation", "get"),
			zap.String("trace", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	var out user.TokenParam

	if err := json.Unmarshal([]byte(res), &out); err != nil {
		err := ose_error.Wrap(err, ose_error.ErrInternalServerError, err.Error(), traceId)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		t.log.Error("error unmarshalling payload",
			zap.String("operation", "get"),
			zap.String("trace", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	return t.toDomain(out)
}

func (t tokenCache) toDomain(param user.TokenParam) (*user.Token, error) {
	return user.ExistingToken(param)
}

func NewTokenCache(client *redis.Client, log logger.Logger, trace tracing.Tracer) user.Cache {
	return tokenCache{
		redis: client,
		log:   log,
		trace: trace,
	}
}

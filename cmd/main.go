package main

import (
	"github.com/ose-micro/authora/internal/api/bus"
	"github.com/ose-micro/authora/internal/api/grpc"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/events"
	"github.com/ose-micro/authora/internal/infrastruture/cache"
	"github.com/ose-micro/authora/internal/infrastruture/repository"
	ose "github.com/ose-micro/core"
	"github.com/ose-micro/core/config"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/tracing"
	ose_jwt "github.com/ose-micro/jwt"
	mongodb "github.com/ose-micro/mongo"
	"github.com/ose-micro/nats"
	"github.com/ose-micro/redis"
	"go.uber.org/fx"
)

func main() {
	ose.New(
		fx.Provide(
			loadConfig,
			mongodb.New,
			nats.New,
			repository.Inject,
			business.Inject,
			app.Inject,
			ose_jwt.NewManager,
			redis.New,
			cache.Inject,
			events.NewEvents,
		),
		fx.Invoke(bus.InvokeConsumers),
		fx.Invoke(grpc.RunGRPCServer),
	).Run()
}

func loadConfig() (config.Service, logger.Config, tracing.Config, timestamp.Config,
	mongodb.Config, nats.Config, grpc.Config, ose_jwt.Config, redis.Config, error) {

	var grpcConfig grpc.Config
	var natsConf nats.Config
	var mongoConfig mongodb.Config
	var jwtConfig ose_jwt.Config
	var redisConfig redis.Config

	conf, err := config.Load(
		config.WithExtension("nats", &natsConf),
		config.WithExtension("mongo", &mongoConfig),
		config.WithExtension("grpc", &grpcConfig),
		config.WithExtension("redis", &redisConfig),
		config.WithExtension("jwt", &jwtConfig),
	)

	if err != nil {
		return config.Service{}, logger.Config{}, tracing.Config{}, timestamp.Config{},
			mongodb.Config{}, nats.Config{}, grpc.Config{}, ose_jwt.Config{}, redis.Config{}, err
	}

	return conf.Core.Service, conf.Core.Service.Logger, conf.Core.Service.Tracer,
		conf.Core.Service.Timestamp, mongoConfig, natsConf, grpcConfig, jwtConfig, redisConfig, nil
}

package main

import (
	"github.com/ose-micro/authora/internal/api/bus"
	"github.com/ose-micro/authora/internal/api/grpc"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/business"
	"github.com/ose-micro/authora/internal/events"
	"github.com/ose-micro/authora/internal/repository"
	ose "github.com/ose-micro/core"
	"github.com/ose-micro/core/config"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/tracing"
	ose_jwt "github.com/ose-micro/jwt"
	mongodb "github.com/ose-micro/mongo"
	"github.com/ose-micro/nats"
	"go.uber.org/fx"
)

func loadConfig() (config.Service, logger.Config, tracing.Config, timestamp.Config,
	mongodb.Config, nats.Config, grpc.Config, ose_jwt.Config, error) {

	var grpcConfig grpc.Config
	var natsConf nats.Config
	var mongoConfig mongodb.Config
	var jwtConfig ose_jwt.Config

	conf, err := config.Load(
		config.WithExtension("nats", &natsConf),
		config.WithExtension("mongo", &mongoConfig),
		config.WithExtension("grpc", &grpcConfig),
		config.WithExtension("jwt", &jwtConfig),
	)

	if err != nil {
		return config.Service{}, logger.Config{}, tracing.Config{}, timestamp.Config{},
			mongodb.Config{}, nats.Config{}, grpc.Config{}, ose_jwt.Config{}, err
	}

	return conf.Core.Service, conf.Core.Service.Logger, conf.Core.Service.Tracer,
		conf.Core.Service.Timestamp, mongoConfig, natsConf, grpcConfig, jwtConfig, nil
}

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
			events.NewEvents,
		),
		fx.Invoke(bus.InvokeConsumers),
		fx.Invoke(grpc.RunGRPCServer),
	).Run()
}

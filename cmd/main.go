package main

import (
	"github.com/ose-micro/authora/internal/api/grpc"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/repository"
	ose "github.com/ose-micro/core"
	"github.com/ose-micro/core/config"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs/bus/nats"
	mongodb "github.com/ose-micro/mongo"
	"go.uber.org/fx"
)

func loadConfig() (config.Service, logger.Config, tracing.Config, timestamp.Config,
	mongodb.Config, nats.Config, grpc.Config, error) {

	var grpcConfig grpc.Config
	var natsConf nats.Config
	var mongoConfig mongodb.Config

	conf, err := config.Load(
		config.WithExtension("bus", &natsConf),
		config.WithExtension("mongo", &mongoConfig),
		config.WithExtension("grpc", &grpcConfig),
	)

	if err != nil {
		return config.Service{}, logger.Config{}, tracing.Config{}, timestamp.Config{},
			mongodb.Config{}, nats.Config{}, grpc.Config{}, err
	}

	return conf.Core.Service, conf.Core.Service.Logger, conf.Core.Service.Tracer,
		conf.Core.Service.Timestamp, mongoConfig, natsConf, grpcConfig, nil
}

func main() {
	ose.New(
		fx.Provide(
			loadConfig,
			mongodb.New,
			nats.New,
			repository.Inject,
			domain.Inject,
			app.Inject,
		),
		fx.Invoke(grpc.RunGRPCServer),
	).Run()
}

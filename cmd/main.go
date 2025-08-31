package main

import (
	"github.com/ose-micro/authora/internal/api/grpc"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/authora/internal/common"
	"github.com/ose-micro/authora/internal/domain"
	"github.com/ose-micro/authora/internal/repository"
	ose "github.com/ose-micro/core"
	"github.com/ose-micro/core/config"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs/bus/nats"
	ose_jwt "github.com/ose-micro/jwt"
	mongodb "github.com/ose-micro/mongo"
	"go.uber.org/fx"
)

func loadConfig() (config.Service, logger.Config, tracing.Config, timestamp.Config,
	mongodb.Config, nats.Config, grpc.Config, ose_jwt.Config, *common.Permissions, error) {

	var grpcConfig grpc.Config
	var natsConf nats.Config
	var mongoConfig mongodb.Config
	var jwtConfig ose_jwt.Config
	var permissionConfig common.Permissions

	conf, err := config.Load(
		config.WithExtension("bus", &natsConf),
		config.WithExtension("mongo", &mongoConfig),
		config.WithExtension("grpc", &grpcConfig),
		config.WithExtension("jwt", &jwtConfig),
		config.WithExtension("permissions", &permissionConfig),
	)

	if err != nil {
		return config.Service{}, logger.Config{}, tracing.Config{}, timestamp.Config{},
			mongodb.Config{}, nats.Config{}, grpc.Config{}, ose_jwt.Config{}, nil, err
	}

	return conf.Core.Service, conf.Core.Service.Logger, conf.Core.Service.Tracer,
		conf.Core.Service.Timestamp, mongoConfig, natsConf, grpcConfig, jwtConfig, &permissionConfig, nil
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
			ose_jwt.NewManager,
		),
		fx.Invoke(grpc.RunGRPCServer),
	).Run()
}

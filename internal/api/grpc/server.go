package grpc

import (
	"context"
	"fmt"
	"net"

	assignmentv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/assignment/v1"
	authv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/auth/v1"
	permissionv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/permission/v1"
	rolev1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/role/v1"
	tenantv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/tenant/v1"
	userv1 "github.com/ose-micro/authora/internal/api/grpc/gen/go/ose/micro/authora/user/v1"
	"github.com/ose-micro/authora/internal/api/grpc/handlers"
	"github.com/ose-micro/authora/internal/app"
	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	osegrpc "github.com/ose-micro/grpc"
	mongodb "github.com/ose-micro/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Port int64 `mapstructure:"port"`
}

func RunGRPCServer(lc fx.Lifecycle, conf Config, log logger.Logger, tracer tracing.Tracer, apps app.Apps, mdb *mongodb.Client,
	bs domain.Bus) (*osegrpc.Server, error) {
	svc, err := osegrpc.New(osegrpc.Params{
		Logger: log,
		Tracer: tracer,
	})
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
				if err != nil {
					log.Fatal("failed to listen", zap.Error(err))
				}

				if err := svc.Serve(lis, func(s *grpc.Server) {
					log.Info(fmt.Sprintf("gRPC server listening on :%d", conf.Port))

					tenantv1.RegisterTenantServiceServer(s, handlers.NewTenant(apps, log, tracer))
					rolev1.RegisterRoleServiceServer(s, handlers.NewRole(apps, log, tracer))
					userv1.RegisterUserServiceServer(s, handlers.NewUser(apps, log, tracer))
					assignmentv1.RegisterAssignmentServiceServer(s, handlers.NewAssignment(apps, log, tracer))
					authv1.RegisterAuthServiceServer(s, handlers.NewAuth(apps, log, tracer))
					permissionv1.RegisterPermissionServiceServer(s, handlers.NewPermission(apps, log, tracer))

				}); err != nil {
					log.Fatal("gRPC server failed", zap.Error(err))
				}

			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := bs.Close()
			if err != nil {
				return err
			}
			log.Info("disconnecting from nats . . .")

			log.Info("disconnecting from mongoDB . . .")
			err = mdb.Close(ctx)
			if err != nil {
				return err
			}

			err = svc.Stop()
			if err != nil {
				return err
			}
			log.Info("gRPC server stopped")

			_, cancel := context.WithCancel(context.Background())
			log.Info("stopping background workers")
			cancel()

			return nil
		},
	})

	return svc, nil
}

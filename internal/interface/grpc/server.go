package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/moriba-cloud/ose-postman/internal/app"
	emailv1 "github.com/moriba-cloud/ose-postman/internal/interface/grpc/gen/go/email/v1"
	templatev1 "github.com/moriba-cloud/ose-postman/internal/interface/grpc/gen/go/template/v1"
	"github.com/moriba-cloud/ose-postman/internal/interface/grpc/handlers"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	oseGrpc "github.com/ose-micro/grpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Port int64 `mapstructure:"port"`
}

func RunGRPCServer(lc fx.Lifecycle, conf Config, log logger.Logger, tracer tracing.Tracer, apps app.Apps) (*oseGrpc.Server, error) {
	svc, err := oseGrpc.New(oseGrpc.Params{
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

					templatev1.RegisterTemplateServiceServer(s, handlers.NewTemplate(apps, log, tracer))
					emailv1.RegisterEmailServiceServer(s, handlers.NewEmail(apps, log, tracer))

				}); err != nil {
					log.Fatal("gRPC server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			svc.Stop()
			log.Info("gRPC server stopped")
			return nil
		},
	})

	return svc, nil
}

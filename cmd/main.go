package main

import (
	ose "github.com/ose-micro/core"
	"github.com/ose-micro/core/config"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/mailer"
	mongodb "github.com/ose-micro/mongo"
	"github.com/ose-micro/nats"
	"github.com/ose-micro/postgres"
	"github.com/ose-micro/postman/internal/api/bus"
	"github.com/ose-micro/postman/internal/api/grpc"
	"github.com/ose-micro/postman/internal/app"
	"github.com/ose-micro/postman/internal/business"
	"github.com/ose-micro/postman/internal/infrastructure/repository"
	"go.uber.org/fx"
)

func main() {
	ose.New(
		fx.Provide(
			loadConfig,
			business.InjectDomain,
			postgres.New,
			mongodb.New,
			mailer.New,
			app.InjectApps,
			nats.New,
			repository.InjectRepository,
		),
		fx.Invoke(bus.InvokeConsumers),
		fx.Invoke(grpc.RunGRPCServer),
	).Run()
}

func loadConfig() (config.Service, logger.Config, tracing.Config, timestamp.Config,
	*postgres.Config, mongodb.Config, nats.Config, grpc.Config, *mailer.Config, error) {
	var grpcConfig grpc.Config
	var natsConf nats.Config
	var postgresConfig postgres.Config
	var mongoConfig mongodb.Config
	var mailerConfig mailer.Config

	conf, err := config.Load(
		config.WithExtension("nats", &natsConf),
		config.WithExtension("postgres", &postgresConfig),
		config.WithExtension("mongo", &mongoConfig),
		config.WithExtension("grpc", &grpcConfig),
		config.WithExtension("mailer", &mailerConfig),
	)

	if err != nil {
		return config.Service{}, logger.Config{}, tracing.Config{}, timestamp.Config{},
			nil, mongodb.Config{}, nats.Config{}, grpc.Config{}, nil, err
	}

	return conf.Core.Service, conf.Core.Service.Logger, conf.Core.Service.Tracer, conf.Core.Service.Timestamp,
		&postgresConfig, mongoConfig, natsConf, grpcConfig, &mailerConfig, nil
}

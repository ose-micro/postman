package main

import (
	"github.com/moriba-cloud/ose-postman/internal/app"
	"github.com/moriba-cloud/ose-postman/internal/domain"
	"github.com/moriba-cloud/ose-postman/internal/events"
	"github.com/moriba-cloud/ose-postman/internal/interface/grpc"
	consumer "github.com/moriba-cloud/ose-postman/internal/interface/bus"
	"github.com/moriba-cloud/ose-postman/internal/repository/read"
	"github.com/moriba-cloud/ose-postman/internal/repository/write"
	ose "github.com/ose-micro/core"
	"github.com/ose-micro/core/config"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/cqrs/bus/rabbitmq"
	"github.com/ose-micro/mailer"
	mongodb "github.com/ose-micro/mongo"
	"github.com/ose-micro/postgres"
	"go.uber.org/fx"
)

func main() {
	app := ose.New(
		fx.Provide(
			loadConfig,
			domain.InjectDomain,
			postgres.New,
			mongodb.New,
			mailer.New,
			app.InjectApps,
			rabbitmq.New,
			events.Inject,
			write.InjectRepository,
			read.InjectRepository,
		),
		fx.Invoke(write.Migrate),
		fx.Invoke(consumer.InvokeConsumers),
		fx.Invoke(grpc.RunGRPCServer),
	)

	app.Run()
}

func loadConfig() (config.Service, logger.Config, tracing.Config, timestamp.Config,
	*postgres.Config, mongodb.Config, rabbitmq.Config, grpc.Config, *mailer.Config, error) {
	var grpcConfig grpc.Config
	var rabbitmqConf rabbitmq.Config
	var postgresConfig postgres.Config
	var mongoConfig mongodb.Config
	var mailerConfig mailer.Config

	conf, err := config.Load(
		config.WithExtension("bus", &rabbitmqConf),
		config.WithExtension("postgres", &postgresConfig),
		config.WithExtension("mongo", &mongoConfig),
		config.WithExtension("grpc", &grpcConfig),
		config.WithExtension("mailer", &mailerConfig),
	)

	if err != nil {
		return config.Service{}, logger.Config{}, tracing.Config{}, timestamp.Config{},
			nil, mongodb.Config{}, rabbitmq.Config{}, grpc.Config{}, nil, err
	}

	return conf.Core.Service, conf.Core.Service.Logger, conf.Core.Service.Tracer, conf.Core.Service.Timestamp, 
	&postgresConfig, mongoConfig, rabbitmqConf, grpcConfig, &mailerConfig, nil
}

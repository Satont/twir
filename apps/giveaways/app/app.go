package app

import (
	"time"

	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/giveaways/internal/grpc_impl"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	twirsentry "github.com/twirapp/twir/libs/sentry"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

var App = fx.Options(
	fx.Provide(
		cfg.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "giveaways"}),
		logger.NewFx(logger.Opts{Service: "giveaways"}),
		uptrace.NewFx("giveaways"),
		func(config cfg.Config) (*gorm.DB, error) {
			db, err := gorm.Open(
				postgres.Open(config.DatabaseUrl), &gorm.Config{
					Logger: gorm_logger.Default.LogMode(gorm_logger.Silent),
				},
			)
			if err != nil {
				return nil, err
			}
			d, _ := db.DB()
			d.SetMaxIdleConns(1)
			d.SetMaxOpenConns(10)
			d.SetConnMaxLifetime(time.Hour)
			return db, nil
		},
		func(config cfg.Config) (*redis.Client, error) {
			redisUrl, err := redis.ParseURL(config.RedisUrl)
			if err != nil {
				return nil, err
			}

			redisClient := redis.NewClient(redisUrl)
			return redisClient, nil
		},
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		buscore.NewNatsBusFx("giveaways"),
	),
	fx.Invoke(
		uptrace.NewFx("giveaways"),
		grpc_impl.New,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)

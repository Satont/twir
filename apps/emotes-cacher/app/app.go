package app

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/emotes-cacher/internal/grpc_impl"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const service = "emotes-cacher"

var App = fx.Module(
	service,
	fx.Provide(
		cfg.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: service}),
		logger.NewFx(logger.Opts{Service: service}),
		uptrace.NewFx("emotes-cacher"),
		func(cfg cfg.Config) (*gorm.DB, error) {
			db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
			if err != nil {
				return nil, err
			}
			d, _ := db.DB()
			d.SetMaxIdleConns(1)
			d.SetMaxOpenConns(10)
			d.SetConnMaxLifetime(time.Hour)

			return db, nil
		},
		func(cfg cfg.Config) (*redis.Client, error) {
			redisUrl, err := redis.ParseURL(cfg.RedisUrl)
			if err != nil {
				return nil, err
			}

			return redis.NewClient(redisUrl), nil
		},
	),
	fx.Invoke(
		uptrace.NewFx("emotes-cacher"),
		grpc_impl.NewEmotesCacher,
		func(l logger.Logger) {
			l.Info("Emotes Cacher started")
		},
	),
)

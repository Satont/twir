package baseapp

import (
	"context"
	"log/slog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/logger/audit"
	auditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	twirsentry "github.com/satont/twir/libs/sentry"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

type Opts struct {
	AppName string

	WithAudit bool
}

func CreateBaseApp(opts Opts) fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewFx,
			twirsentry.NewFx(twirsentry.NewFxOpts{Service: opts.AppName}),
			uptrace.NewFx(opts.AppName),
			newRedis,
			newGorm,
			newPgxPool,
			buscore.NewNatsBusFx(opts.AppName),
			fx.Annotate(
				auditlogs.NewBusPubSubFx,
				fx.As(new(auditlogs.PubSub)),
			),
			fx.Annotate(
				audit.NewGormFx(),
				fx.As(new(slog.Handler)),
				fx.ResultTags(`group:"slog-handlers"`),
			),
			fx.Annotate(
				audit.NewPubsubFx(),
				fx.As(new(slog.Handler)),
				fx.ResultTags(`group:"slog-handlers"`),
			),
			logger.NewFx(
				logger.Opts{
					Service: opts.AppName,
					Level:   slog.LevelInfo,
				},
			),
		),
		fx.Invoke(uptrace.NewFx(opts.AppName)),
	)
}

func newRedis(cfg config.Config) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(redisOpts)

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return nil, err
	}

	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		return nil, err
	}

	return redisClient, nil
}

func newPgxPool(cfg config.Config) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	connConfig.MaxConnLifetime = time.Hour
	connConfig.MaxConnIdleTime = 5 * time.Minute
	connConfig.MaxConns = 100
	connConfig.MinConns = 1
	connConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		connConfig,
	)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	bus_listener "github.com/satont/twir/apps/bots/internal/bus-listener"
	"github.com/satont/twir/apps/bots/internal/gorm"
	"github.com/satont/twir/apps/bots/internal/messagehandler"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	stream_handlers "github.com/satont/twir/apps/bots/internal/stream-handlers"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/giveaways"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"bots",
	fx.Provide(
		cfg.NewFx,
		tlds.New,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "bots"}),
		logger.NewFx(logger.Opts{Service: "bots"}),
		gorm.New,
		uptrace.NewFx("bots"),
		buscore.NewNatsBusFx("bots"),
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		func(config cfg.Config) events.EventsClient {
			return clients.NewEvents(config.AppEnv)
		},
		func(config cfg.Config) parser.ParserClient {
			return clients.NewParser(config.AppEnv)
		},
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		func(config cfg.Config) giveaways.GiveawaysClient {
			return clients.NewGiveaways(config.AppEnv)
		},
		func(config cfg.Config) (*redis.Client, error) {
			redisOpts, err := redis.ParseURL(config.RedisUrl)
			if err != nil {
				return nil, err
			}

			return redis.NewClient(redisOpts), nil
		},
		twitchactions.New,
		moderationhelpers.New,
		messagehandler.New,
	),
	fx.Invoke(
		uptrace.NewFx("bots"),
		func(config cfg.Config) {
			if config.AppEnv != "development" {
				http.Handle("/metrics", promhttp.Handler())
				go http.ListenAndServe("0.0.0.0:3000", nil)
			}
		},
		stream_handlers.New,
		bus_listener.New,
		func(l logger.Logger) {
			l.Info("Bots started")
		},
	),
)

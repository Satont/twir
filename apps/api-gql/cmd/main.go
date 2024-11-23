package main

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/gql"
	community_searcher "github.com/twirapp/twir/apps/api-gql/internal/gql/community-searcher"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/resolvers"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/gql/twir-stats"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/routes/auth"
	pubclicroutes "github.com/twirapp/twir/apps/api-gql/internal/routes/public"
	"github.com/twirapp/twir/apps/api-gql/internal/routes/webhooks"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	"github.com/twirapp/twir/libs/baseapp"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	keywordscacher "github.com/twirapp/twir/libs/cache/keywords"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "api-gql",
			},
		),
		fx.Provide(
			auth.NewSessions,
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			func(config cfg.Config) events.EventsClient {
				return clients.NewEvents(config.AppEnv)
			},
			minio.New,
			twitchcache.New,
			commandscache.New,
			keywordscacher.New,
			fx.Annotate(
				wsrouter.NewNatsSubscription,
				fx.As(new(wsrouter.WsRouter)),
			),
			twir_stats.New,
			resolvers.New,
			directives.New,
			httpserver.New,
			gql.New,
			community_searcher.NewCommunitySearcher,
		),
		fx.Invoke(
			pubclicroutes.New,
			webhooks.New,
			authroutes.New,
		),
	).Run()
}

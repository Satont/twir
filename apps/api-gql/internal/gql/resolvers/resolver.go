package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/nicklaw5/helix/v2"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/resolvers/integrations"
	subscriptions_store "github.com/twirapp/twir/apps/api-gql/internal/gql/subscriptions-store"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	config                      config.Config
	sessions                    *sessions.Sessions
	gorm                        *gorm.DB
	twitchClient                *helix.Client
	cachedTwitchClient          *twitchcahe.CachedTwitchClient
	minioClient                 *minio.Client
	subscriptionsStore          *subscriptions_store.SubscriptionsStore
	integrationsLinksResolver   *integrations.LinksResolver
	integrationsDataFetcher     *integrations.DataFetcher
	integrationsPostCodeHandler *integrations.PostCodeHandler
}

type Opts struct {
	fx.In

	Sessions                    *sessions.Sessions
	Gorm                        *gorm.DB
	Config                      config.Config
	TokensGrpc                  tokens.TokensClient
	CachedTwitchClient          *twitchcahe.CachedTwitchClient
	Minio                       *minio.Client
	SubscriptionsStore          *subscriptions_store.SubscriptionsStore
	IntegrationsLinksResolver   *integrations.LinksResolver
	IntegrationsDataFetcher     *integrations.DataFetcher
	IntegrationsPostCodeHandler *integrations.PostCodeHandler
}

func New(opts Opts) (*Resolver, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.TokensGrpc)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		config:                      opts.Config,
		sessions:                    opts.Sessions,
		gorm:                        opts.Gorm,
		twitchClient:                twitchClient,
		cachedTwitchClient:          opts.CachedTwitchClient,
		minioClient:                 opts.Minio,
		subscriptionsStore:          opts.SubscriptionsStore,
		integrationsLinksResolver:   opts.IntegrationsLinksResolver,
		integrationsDataFetcher:     opts.IntegrationsDataFetcher,
		integrationsPostCodeHandler: opts.IntegrationsPostCodeHandler,
	}, nil
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(
	ctx *graphql.OperationContext,
	fields []graphql.CollectedField,
	prefix string,
) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(
			preloads,
			GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...,
		)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}

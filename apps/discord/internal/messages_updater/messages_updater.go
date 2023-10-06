package messages_updater

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Store   *sended_messages_store.SendedMessagesStore
	Logger  logger.Logger
	LC      fx.Lifecycle
	Config  cfg.Config
	DB      *gorm.DB
	Discord *discord_go.Discord

	TokensGrpc tokens.TokensClient
}

func New(opts Opts) *MessagesUpdater {
	updater := &MessagesUpdater{
		store:      opts.Store,
		logger:     opts.Logger.WithComponent("messages_updater"),
		config:     opts.Config,
		db:         opts.DB,
		discord:    opts.Discord,
		tokensGrpc: opts.TokensGrpc,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go updater.poll(ctx)
				updater.logger.Info("Messages updater is running")

				return nil
			},
		},
	)

	return updater
}

type MessagesUpdater struct {
	store   *sended_messages_store.SendedMessagesStore
	logger  logger.Logger
	config  cfg.Config
	db      *gorm.DB
	discord *discord_go.Discord

	tokensGrpc tokens.TokensClient
}

func (c *MessagesUpdater) poll(ctx context.Context) {
	ticker := time.NewTicker(
		lo.If(
			c.config.AppEnv != "production",
			10*time.Second,
		).Else(5 * time.Minute),
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.process(ctx)
		}
	}
}
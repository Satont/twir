package messagehandler

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"runtime"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/giveaways"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger            logger.Logger
	Gorm              *gorm.DB
	Redis             *redis.Client
	TwitchActions     *twitchactions.TwitchActions
	ParserGrpc        parser.ParserClient
	WebsocketsGrpc    websockets.WebsocketClient
	EventsGrpc        events.EventsClient
	ModerationHelpers *moderationhelpers.ModerationHelpers
	Config            cfg.Config
	Bus               *buscore.Bus
	GiveawaysGrpc     giveaways.GiveawaysClient
}

type MessageHandler struct {
	logger            logger.Logger
	gorm              *gorm.DB
	redis             *redis.Client
	twitchActions     *twitchactions.TwitchActions
	parserGrpc        parser.ParserClient
	websocketsGrpc    websockets.WebsocketClient
	eventsGrpc        events.EventsClient
	moderationHelpers *moderationhelpers.ModerationHelpers
	config            cfg.Config
	bus               *buscore.Bus
	giveawaysGrpc     giveaways.GiveawaysClient
}

func New(opts Opts) *MessageHandler {
	handler := &MessageHandler{
		logger:            opts.Logger,
		gorm:              opts.Gorm,
		redis:             opts.Redis,
		twitchActions:     opts.TwitchActions,
		parserGrpc:        opts.ParserGrpc,
		websocketsGrpc:    opts.WebsocketsGrpc,
		eventsGrpc:        opts.EventsGrpc,
		moderationHelpers: opts.ModerationHelpers,
		config:            opts.Config,
		bus:               opts.Bus,
		giveawaysGrpc:     opts.GiveawaysGrpc,
	}

	return handler
}

type handleMessage struct {
	twitch.TwitchChatMessage
	DbChannel *model.Channels
	DbStream  *model.ChannelsStreams
	DbUser    *model.Users
}

var handlersForExecute = []func(
	c *MessageHandler,
	ctx context.Context,
	msg handleMessage,
) error{
	(*MessageHandler).handleIncrementStreamMessages,
	(*MessageHandler).handleGreetings,
	(*MessageHandler).handleKeywords,
	(*MessageHandler).handleEmotesUsages,
	(*MessageHandler).handleStoreMessage,
	(*MessageHandler).handleTts,
	(*MessageHandler).handleRemoveLurker,
	(*MessageHandler).handleModeration,
	(*MessageHandler).handleFirstStreamUserJoin,
	(*MessageHandler).handleGiveaways,
}

func (c *MessageHandler) Handle(ctx context.Context, req twitch.TwitchChatMessage) error {
	msg := handleMessage{
		TwitchChatMessage: req,
	}

	errwg, errWgCtx := errgroup.WithContext(context.TODO())

	errwg.Go(
		func() error {
			stream := &model.ChannelsStreams{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				`"userId" = ?`,
				req.BroadcasterUserId,
			).First(stream).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if stream.ID == "" {
				msg.DbStream = nil
			} else {
				msg.DbStream = stream
			}
			return nil
		},
	)

	errwg.Go(
		func() error {
			dbChannel := &model.Channels{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				"id = ?",
				req.BroadcasterUserId,
			).First(dbChannel).
				Error; err != nil {
				return err
			}
			msg.DbChannel = dbChannel
			return nil
		},
	)

	if err := errwg.Wait(); err != nil {
		return err
	}

	if !msg.DbChannel.IsEnabled {
		return nil
	}

	dbUser, err := c.ensureUser(ctx, msg)
	if err != nil {
		return err
	}
	msg.DbUser = dbUser

	if req.ChatterUserId == msg.DbChannel.BotID && c.config.AppEnv == "production" {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(handlersForExecute))

	// TODO: i dont know why grpc context canceling before this function finished
	funcsCtx := context.Background()

	for _, f := range handlersForExecute {
		f := f

		go func() {
			defer wg.Done()
			if err := f(c, funcsCtx, msg); err != nil {
				c.logger.Error(
					"error when executing message handler function",
					slog.Any("err", err),
					slog.String(
						"functionName",
						runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
					),
				)
			}
		}()
	}

	wg.Wait()

	return nil
}

package impl_deps

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/giveaways"
	"github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"gorm.io/gorm"
)

type Grpc struct {
	Tokens       tokens.TokensClient
	Integrations integrations.IntegrationsClient
	Parser       parser.ParserClient
	Websockets   websockets.WebsocketClient
	Discord      discord.DiscordClient
	Giveaways    giveaways.GiveawaysClient
}

type Deps struct {
	Config         cfg.Config
	Redis          *redis.Client
	Db             *gorm.DB
	Grpc           *Grpc
	SessionManager *scs.SessionManager
	Logger         logger.Logger
	Bus            *buscore.Bus
}

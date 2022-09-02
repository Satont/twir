package variables

import (
	"regexp"
	"strings"
	"sync"
	"tsuwari/parser/internal/config/twitch"
	types "tsuwari/parser/internal/types"
	emotes7tv "tsuwari/parser/internal/variables/emotes/7tv"
	emotesbttv "tsuwari/parser/internal/variables/emotes/bttv"
	emotesffz "tsuwari/parser/internal/variables/emotes/ffz"
	"tsuwari/parser/internal/variables/random"
	sender "tsuwari/parser/internal/variables/sender"
	streamcategory "tsuwari/parser/internal/variables/stream/category"
	streamtitle "tsuwari/parser/internal/variables/stream/title"
	streamuptime "tsuwari/parser/internal/variables/stream/uptime"
	streamviewers "tsuwari/parser/internal/variables/stream/viewers"
	userage "tsuwari/parser/internal/variables/user/age"
	userfollowage "tsuwari/parser/internal/variables/user/followage"
	usermessages "tsuwari/parser/internal/variables/user/messages"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type Variables struct {
	Store  map[string]types.Variable
	Redis  *redis.Client
	Twitch *twitch.Twitch
	Db     *gorm.DB
}

var Regexp = regexp.MustCompile(`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`)

func New(redis *redis.Client, twitchApi *twitch.Twitch, db *gorm.DB) Variables {
	ctx := Variables{
		Store:  make(map[string]types.Variable),
		Redis:  redis,
		Twitch: twitchApi,
		Db:     db,
	}

	ctx.Store[random.Name] = types.Variable{
		Name:    random.Name,
		Handler: random.Handler,
	}
	ctx.Store[sender.Name] = types.Variable{
		Name:    sender.Name,
		Handler: sender.Handler,
	}
	ctx.Store[streamuptime.Name] = types.Variable{
		Name:    streamuptime.Name,
		Handler: streamuptime.Handler,
	}
	ctx.Store[streamtitle.Name] = types.Variable{
		Name:    streamtitle.Name,
		Handler: streamtitle.Handler,
	}
	ctx.Store[streamviewers.Name] = types.Variable{
		Name:    streamviewers.Name,
		Handler: streamviewers.Handler,
	}
	ctx.Store[streamcategory.Name] = types.Variable{
		Name:    streamcategory.Name,
		Handler: streamcategory.Handler,
	}
	ctx.Store[emotesffz.Name] = types.Variable{
		Name:    emotesffz.Name,
		Handler: emotesffz.Handler,
	}
	ctx.Store[emotes7tv.Name] = types.Variable{
		Name:    emotes7tv.Name,
		Handler: emotes7tv.Handler,
	}
	ctx.Store[emotesbttv.Name] = types.Variable{
		Name:    emotesbttv.Name,
		Handler: emotesbttv.Handler,
	}
	ctx.Store[usermessages.Name] = types.Variable{
		Name:    usermessages.Name,
		Handler: usermessages.Handler,
	}
	ctx.Store[userfollowage.Name] = types.Variable{
		Name:    userfollowage.Name,
		Handler: userfollowage.Handler,
	}
	ctx.Store[userage.Name] = types.Variable{
		Name:    userage.Name,
		Handler: userage.Handler,
	}

	return ctx
}

func (c Variables) ParseInput(cache *variablescache.VariablesCacheService, input string) string {
	wg := sync.WaitGroup{}

	for _, s := range Regexp.FindAllString(input, len(input)) {
		wg.Add(1)
		v := Regexp.FindStringSubmatch(s)
		all := v[1]
		params := v[3]

		if val, ok := c.Store[all]; ok {
			go func(s string) {
				defer wg.Done()
				res, err := val.Handler(cache, types.VariableHandlerParams{
					Key:    all,
					Params: &params,
				})

				if err == nil {
					input = strings.ReplaceAll(input, s, res.Result)
				}
			}(s)
		} else {
			wg.Done()
		}
	}

	wg.Wait()
	return input
}

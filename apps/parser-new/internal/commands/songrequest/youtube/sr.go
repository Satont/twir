package sr_youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/libs/grpc/generated/ytsr"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/valyala/fasttemplate"
	"go.uber.org/zap"

	"github.com/satont/tsuwari/libs/grpc/generated/websockets"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/nicklaw5/helix/v2"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/samber/lo"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
)

type ReqError struct {
	Title string
	Error string
}

var SrCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr",
		Description: null.StringFrom("Song requests from youtube"),
		Module:      "SONGS",
		IsReply:     true,
		Visible:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		if parseCtx.Text == nil {
			result.Result = append(result.Result, "You should provide text for song request")
			return result
		}

		moduleSettings := &model.ChannelModulesSettings{}
		parsedSettings := &youtube.YouTubeSettings{}
		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "type" = ?`, parseCtx.Channel.ID, "youtube_song_requests").
			First(moduleSettings).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "internal error")
			return result
		}
		if moduleSettings.ID != "" {
			err = json.Unmarshal(moduleSettings.Settings, parsedSettings)
			if err != nil {
				zap.S().Error(err)
				result.Result = append(result.Result, "internal error")
				return result
			}

			if !*parsedSettings.Enabled {
				result.Result = append(result.Result, parsedSettings.Translations.NotEnabled)
				return result
			}
		} else {
			result.Result = append(result.Result, "Song requests not enabled")
			return result
		}

		req, err := parseCtx.Services.GrpcClients.Ytsr.Search(context.Background(), &ytsr.SearchRequest{Search: *parseCtx.Text})
		if err != nil {
			zap.S().Error(err)
			return result
		}
		if len(req.Songs) == 0 {
			result.Result = append(result.Result, parsedSettings.Translations.Song.NotFound)
			return result
		}

		latestSong := &model.RequestedSong{}

		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, parseCtx.Channel.ID).
			Order(`"createdAt" desc`).
			Find(&latestSong).Error

		requested := make([]*model.RequestedSong, 0, len(req.Songs))
		errors := make([]*ReqError, 0, len(req.Songs))

		var currentQueueCount int64
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, parseCtx.Channel.ID).
			Model(&model.RequestedSong{}).
			Count(&currentQueueCount).
			Error

		if err != nil {
			zap.S().Error(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		for i, song := range req.Songs {
			err = validate(
				ctx,
				parseCtx.Services,
				parseCtx.Channel.ID,
				parseCtx.Sender.ID,
				parsedSettings,
				song,
			)

			if err != nil {
				errors = append(errors, &ReqError{
					Title: song.Title,
					Error: err.Error(),
				})
			} else {
				model := &model.RequestedSong{
					ID:                   uuid.NewV4().String(),
					ChannelID:            parseCtx.Channel.ID,
					OrderedById:          parseCtx.Sender.ID,
					OrderedByName:        parseCtx.Sender.Name,
					OrderedByDisplayName: null.StringFrom(parseCtx.Sender.DisplayName),
					VideoID:              song.Id,
					Title:                song.Title,
					Duration:             int32(song.Duration),
					CreatedAt:            time.Now().UTC(),
					QueuePosition:        int(currentQueueCount) + (i + 1),
				}

				err = parseCtx.Services.Gorm.WithContext(ctx).Create(model).Error
				if err == nil {
					requested = append(requested, model)
				}
			}
		}

		if len(requested) > 0 {
			requestedMapped := lo.Map(requested, func(item *model.RequestedSong, _ int) string {
				return fmt.Sprintf("%s (#%v)", item.Title, item.QueuePosition)
			})

			result.Result = append(result.Result, "✅ "+strings.Join(requestedMapped, " · "))
		}

		if len(errors) > 0 {
			errorsMapped := lo.Map(errors, func(item *ReqError, _ int) string {
				return item.Title + " - " + item.Error
			})
			result.Result = append(result.Result, "❌"+strings.Join(errorsMapped, " · "))
		}

		for _, song := range requested {
			parseCtx.Services.GrpcClients.WebSockets.YoutubeAddSongToQueue(
				context.Background(),
				&websockets.YoutubeAddSongToQueueRequest{
					ChannelId: parseCtx.Channel.ID,
					EntityId:  song.ID,
				},
			)
		}

		return result
	},
}

func validate(
	ctx context.Context,
	services *services.Services,
	channelId, userId string,
	settings *youtube.YouTubeSettings,
	song *ytsr.Song,
) error {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, *services.Config, services.GrpcClients.Tokens)
	if err != nil {
		return err
	}

	alreadyRequestedSong := &model.RequestedSong{}
	services.Gorm.WithContext(ctx).Where(`"videoId" = ? AND "deletedAt" IS NULL AND "channelId" = ?`, song.Id, channelId).
		Find(&alreadyRequestedSong)

	if alreadyRequestedSong.ID != "" {
		return errors.New(settings.Translations.Song.AlreadyInQueue)
	}

	if song.IsLive {
		return errors.New(settings.Translations.Song.Live)
	}

	if len(settings.DenyList.Users) > 0 {
		_, isUserDenied := lo.Find(
			settings.DenyList.Users,
			func(u youtube.YouTubeDenySettingsUsers) bool {
				return u.UserID == userId
			},
		)

		if isUserDenied {
			return errors.New(settings.Translations.User.Denied)
		}
	}

	if len(settings.DenyList.Channels) > 0 && song.Author != nil {
		_, isChannelBlacklisted := lo.Find(
			settings.DenyList.Channels,
			func(u youtube.YouTubeDenySettingsChannels) bool {
				return u.ID == song.Author.ChannelId
			},
		)

		if isChannelBlacklisted {
			return errors.New(settings.Translations.Channel.Denied)
		}
	}

	if len(settings.DenyList.Songs) > 0 {
		_, isSongBlackListed := lo.Find(
			settings.DenyList.Songs,
			func(u youtube.YouTubeDenySettingsSongs) bool {
				return u.ID == song.Id
			},
		)

		if isSongBlackListed {
			return errors.New(settings.Translations.Song.Denied)
		}
	}

	if *settings.AcceptOnlyWhenOnline {
		stream := &model.ChannelsStreams{}
		services.Gorm.WithContext(ctx).Where(`"userId" = ?`, channelId).First(stream)
		if stream.ID == "" {
			return errors.New(settings.Translations.AcceptOnlineWhenOnline)
		}
	}

	if settings.MaxRequests != 0 {
		var count int64
		services.Gorm.WithContext(ctx).Model(&model.RequestedSong{}).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
			Count(&count)
		if count >= int64(settings.MaxRequests) {
			message := fasttemplate.ExecuteString(
				settings.Translations.Song.MaximumOrdered,
				"{{", "}}",
				map[string]interface{}{
					"maximum": strconv.Itoa(settings.MaxRequests),
				},
			)
			return errors.New(message)
		}
	}

	if settings.Song.MinViews != 0 && int(song.Views) < settings.Song.MinViews {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MinViews,
			"{{", "}}",
			map[string]interface{}{
				"songTitle":   song.Title,
				"songId":      song.Id,
				"songViews":   strconv.Itoa(int(song.Views)),
				"neededViews": strconv.Itoa(settings.Song.MinViews),
			},
		)
		return errors.New(message)
	}

	songDuration := time.Duration(song.Duration) * time.Millisecond
	if settings.Song.MaxLength != 0 && int(math.Round(songDuration.Minutes())) > settings.Song.MaxLength {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MaxLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.Id,
				"songViews": strconv.Itoa(int(song.Views)),
				"maxLength": strconv.Itoa(settings.Song.MaxLength),
			},
		)
		return errors.New(message)
	}

	if settings.Song.MinLength != 0 && int(math.Round(songDuration.Minutes())) < settings.Song.MinLength {
		message := fasttemplate.ExecuteString(
			settings.Translations.Song.MinLength,
			"{{", "}}",
			map[string]interface{}{
				"songTitle": song.Title,
				"songId":    song.Id,
				"songViews": strconv.Itoa(int(song.Views)),
				"minLength": strconv.Itoa(settings.Song.MinLength),
			},
		)
		return errors.New(message)
	}

	// TODO: check categories

	if settings.User.MaxRequests != 0 {
		var count int64
		services.Gorm.WithContext(ctx).
			Model(&model.RequestedSong{}).
			Where(`"orderedById" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, userId, channelId).
			Count(&count)
		if count >= int64(settings.User.MaxRequests) {
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MaxRequests,
				"{{", "}}",
				map[string]interface{}{
					"count":   strconv.Itoa(settings.User.MaxRequests),
					"maximum": strconv.Itoa(settings.User.MaxRequests),
				},
			)
			return errors.New(message)
		}
	}

	if settings.User.MinMessages != 0 || settings.User.MinWatchTime != 0 {
		user := &model.Users{}
		services.Gorm.WithContext(ctx).Where("id = ?", userId).Preload("Stats").First(&user)
		if user.ID == "" {
			return errors.New(
				"there is restrictions on user, but i cannot find you in db, sorry. :(",
			)
		}

		if settings.User.MinMessages != 0 &&
			user.Stats.Messages < int32(settings.User.MinMessages) {
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MinMessages,
				"{{", "}}",
				map[string]interface{}{
					"minMessages":  strconv.Itoa(settings.User.MinMessages),
					"userMessages": strconv.Itoa(int(user.Stats.Messages)),
				},
			)
			return errors.New(message)
		}

		watchedInMinutes := time.Duration(user.Stats.Watched) * time.Millisecond
		if settings.User.MinWatchTime != 0 &&
			int64(watchedInMinutes.Minutes()) < settings.User.MinWatchTime {
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MinWatched,
				"{{", "}}",
				map[string]interface{}{
					"minWatched":  strconv.Itoa(int(settings.User.MinWatchTime)),
					"userWatched": strconv.Itoa(int(watchedInMinutes)),
				},
			)
			return errors.New(message)
		}
	}

	if settings.User.MinFollowTime != 0 {
		neededDuration := time.Minute * time.Duration(settings.User.MinFollowTime)
		followReq, err := twitchClient.GetUsersFollows(&helix.UsersFollowsParams{
			FromID: userId,
			ToID:   channelId,
		})
		if err != nil {
			return errors.New("Internal error when checking follow")
		}
		if followReq.Data.Total == 0 {
			return errors.New("For request song you need to be a followed")
		}

		followDuration := time.Since(followReq.Data.Follows[0].FollowedAt)
		if followDuration.Minutes() < neededDuration.Minutes() {
			message := fasttemplate.ExecuteString(
				settings.Translations.User.MinFollow,
				"{{", "}}",
				map[string]interface{}{
					"minFollow":  strconv.Itoa(settings.User.MinFollowTime),
					"userFollow": strconv.Itoa(int(followDuration.Minutes())),
				},
			)
			return errors.New(message)
		}
	}

	return nil
}

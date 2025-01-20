package song_request

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/ytsr"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm     *gorm.DB
	TwirBus  *buscore.Bus
	YtsrGrpc ytsr.YtsrClient
	Logger   logger.Logger
}

func New(opts Opts) *SongRequest {
	return &SongRequest{
		gorm:     opts.Gorm,
		twirBus:  opts.TwirBus,
		ytsrGrpc: opts.YtsrGrpc,
		logger:   opts.Logger,
	}
}

type SongRequest struct {
	gorm     *gorm.DB
	twirBus  *buscore.Bus
	ytsrGrpc ytsr.YtsrClient
	logger   logger.Logger
}

type ProcessFromDonationInput struct {
	Text      string
	ChannelID string
}

func (c *SongRequest) ProcessFromDonation(
	ctx context.Context,
	input ProcessFromDonationInput,
) error {
	srSettings := model.ChannelSongRequestsSettings{}
	if err := c.gorm.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			input.ChannelID,
		).
		First(&srSettings).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("cannot get song request settings: %w", err)
	}

	srCommand := model.ChannelsCommands{}
	if err := c.gorm.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "defaultName" = ?`,
			input.ChannelID,
			"sr",
		).
		First(&srCommand).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("cannot get song request command: %w", err)
	}

	if !srCommand.Enabled || !srSettings.Enabled || !srSettings.TakeSongFromDonationMessage {
		return nil
	}

	ytsrResult, err := c.ytsrGrpc.Search(
		ctx, &ytsr.SearchRequest{
			Search:    input.Text,
			OnlyLinks: true,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot search for ytsrResult: %w", err)
	}

	for _, song := range ytsrResult.GetSongs() {
		err := c.twirBus.Parser.ProcessMessageAsCommand.Publish(
			twitch.TwitchChatMessage{
				ID:                   "",
				BroadcasterUserId:    input.ChannelID,
				BroadcasterUserName:  "",
				BroadcasterUserLogin: "",
				ChatterUserId:        input.ChannelID,
				ChatterUserName:      "",
				ChatterUserLogin:     "",
				MessageId:            "",
				Message: &twitch.ChatMessageMessage{
					Text: fmt.Sprintf(
						"!%s https://youtu.be/%s",
						srCommand.Name,
						song.GetId(),
					),
					Fragments: nil,
				},
				Color: "",
				Badges: []twitch.ChatMessageBadge{
					{
						Id:    "BROADCASTER",
						SetId: "BROADCASTER",
						Info:  "BROADCASTER",
					},
				},
				MessageType:                 "",
				Cheer:                       nil,
				Reply:                       nil,
				ChannelPointsCustomRewardId: "",
			},
		)

		if err != nil {
			c.logger.Error("cannot publish process message", slog.Any("err", err))
		}
	}

	return nil
}
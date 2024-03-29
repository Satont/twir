package messagehandler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/grpc/giveaways"
)

func (c *MessageHandler) handleGiveaways(ctx context.Context, msg handleMessage) error {
	if msg.DbStream == nil {
		return nil
	}

	go func() {
		_, err := c.giveawaysGrpc.TryProcessParticipant(
			ctx,
			&giveaways.TryProcessParticipantRequest{
				ChannelId:   msg.DbChannel.ID,
				UserId:      msg.DbUser.ID,
				MessageText: msg.Message.Text,
				DisplayName: msg.ChatterUserName,
			},
		)
		if err != nil {
			c.logger.Error(
				"cannot process participant",
				slog.String("channelId", msg.BroadcasterUserId),
				slog.Any("err", err),
			)
		}
	}()

	return nil
}

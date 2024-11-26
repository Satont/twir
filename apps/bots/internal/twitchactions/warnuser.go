package twitchactions

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

type WarnUserOpts struct {
	BroadcasterID string
	ModeratorID   string
	UserID        string
	Reason        string
}

func (c *TwitchActions) WarnUser(ctx context.Context, opts WarnUserOpts) error {
	twitchClient, err := twitch.NewBotClientWithContext(ctx, opts.ModeratorID, c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	resp, err := twitchClient.SendModeratorWarnMessage(
		&helix.SendModeratorWarnChatMessageParams{
			BroadcasterID: opts.BroadcasterID,
			ModeratorID:   opts.ModeratorID,
			UserID:        opts.UserID,
			Reason:        opts.Reason,
		},
	)
	if err != nil {
		return err
	}

	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}

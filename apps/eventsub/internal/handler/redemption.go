package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/websockets"

	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelPointsRewardRedemptionAdd(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) {
	c.logger.Info(
		"channel points reward redemption add",
		slog.String("reward", event.Reward.Title),
		slog.String("userName", event.UserLogin),
		slog.String("userId", event.UserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("channelId", event.BroadcasterUserID),
	)

	err := c.gorm.Create(
		&model.ChannelRedemption{
			ID:           uuid.MustParse(event.ID),
			ChannelID:    event.BroadcasterUserID,
			UserID:       event.UserID,
			RewardID:     uuid.MustParse(event.Reward.ID),
			RewardTitle:  event.Reward.Title,
			RewardPrompt: null.StringFrom(event.UserInput),
			RewardCost:   event.Reward.Cost,
			RedeemedAt:   time.Now().UTC(),
		},
	).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	err = c.gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			UserID:    event.UserID,
			Type:      model.ChannelEventListItemTypeRedemptionCreated,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				RedemptionInput:           event.UserInput,
				RedemptionTitle:           event.Reward.Title,
				RedemptionUserName:        event.UserLogin,
				RedemptionUserDisplayName: event.UserName,
				RedemptionCost:            strconv.Itoa(event.Reward.Cost),
			},
		},
	).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	// fire event to events microsevice
	_, err = c.eventsGrpc.RedemptionCreated(
		context.Background(),
		&events.RedemptionCreatedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Id:              event.Reward.ID,
			RewardName:      event.Reward.Title,
			RewardCost:      strconv.Itoa(event.Reward.Cost),
			Input:           lo.If(event.UserInput != "", &event.UserInput).Else(nil),
			UserId:          event.UserID,
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	// update user spend points
	go func() {
		e := c.countUserChannelPoints(event.UserID, event.BroadcasterUserID, event.Reward.Cost)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("err", e))
		}
	}()

	// youtube song requests

	go func() {
		e := c.handleYoutubeSongRequests(event)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("e", err))
		}
	}()

	go func() {
		e := c.handleAlerts(event)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("e", err))
		}
	}()

	go func() {
		e := c.handleRewardsSevenTvEmote(event)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("err", e))
		}
	}()

}

func (c *Handler) handleChannelPointsRewardRedemptionUpdate(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionUpdate,
) {
	if event.Status != "CANCELED" {
		return
	}

	userStats := &model.UsersStats{}
	err := c.gorm.Where(`"userId" = ?`, event.UserID).Find(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
	if userStats.ID == "" {
		return
	}
	userStats.UsedChannelPoints -= int64(event.Reward.Cost)
	err = c.gorm.Save(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
}

func (c *Handler) countUserChannelPoints(userId, channelId string, count int) error {
	user := &model.Users{}
	err := c.gorm.
		Where("id = ?", userId).
		Preload("Stats", `"channelId" = ?`, channelId).
		First(user).Error
	if err != nil {
		return err
	}

	if user.ID == "" {
		user = &model.Users{
			ID:         "",
			TokenID:    sql.NullString{},
			IsTester:   false,
			IsBotAdmin: false,
			ApiKey:     uuid.New().String(),
			Stats: &model.UsersStats{
				ID:                uuid.New().String(),
				UserID:            userId,
				ChannelID:         channelId,
				Messages:          0,
				Watched:           0,
				UsedChannelPoints: int64(count),
				Emotes:            0,
			},
		}

		err = c.gorm.Error
		if err != nil {
			return err
		}
	}

	if user.Stats != nil {
		user.Stats.UsedChannelPoints += int64(count)
		err = c.gorm.Save(user.Stats).Error
		if err != nil {
			return err
		}
	} else {
		user.Stats = &model.UsersStats{
			ID:                uuid.New().String(),
			UserID:            userId,
			ChannelID:         channelId,
			Messages:          0,
			Watched:           0,
			UsedChannelPoints: int64(count),
			Emotes:            0,
		}
		err = c.gorm.Create(user.Stats).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Handler) handleYoutubeSongRequests(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	if event.UserInput == "" {
		return nil
	}

	entity := &model.ChannelSongRequestsSettings{}
	err := c.gorm.
		Where(`"channel_id" = ?`, event.BroadcasterUserID).
		Find(entity).
		Error
	if err != nil {
		return err
	}
	if entity.ID == "" {
		return nil
	}

	if !entity.Enabled || event.Reward.ID != entity.ChannelPointsRewardID.String {
		return nil
	}

	command := &model.ChannelsCommands{}
	err = c.gorm.
		Where(`"channelId" = ? AND "defaultName" = ?`, event.BroadcasterUserID, "sr").
		Find(command).Error
	if err != nil {
		return err
	}
	if command.ID == "" {
		c.logger.Warn("no command sr", slog.String("channelId", event.BroadcasterUserID))
		return nil
	}

	res, err := c.bus.Parser.GetCommandResponse.Request(
		context.Background(), twitch.TwitchChatMessage{
			BroadcasterUserId:    event.BroadcasterUserID,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ChatterUserId:        event.UserID,
			ChatterUserName:      event.UserName,
			ChatterUserLogin:     event.UserLogin,
			MessageId:            event.ID,
			Message: &twitch.ChatMessageMessage{
				Text: fmt.Sprintf("!%s %s", command.Name, event.UserInput),
			},
		},
	)
	if err != nil {
		return err
	}

	if len(res.Data.Responses) == 0 {
		return nil
	}

	for _, response := range res.Data.Responses {
		c.bus.Bots.SendMessage.Publish(
			bots.SendMessageRequest{
				ChannelId:      event.BroadcasterUserID,
				ChannelName:    &event.BroadcasterUserLogin,
				Message:        fmt.Sprintf("@%s %s", event.UserLogin, response),
				SkipRateLimits: true,
			},
		)
	}

	return nil
}

func (c *Handler) handleAlerts(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	alert := model.ChannelAlert{}

	if err := c.gorm.Where(
		"channel_id = ? AND reward_ids && ?",
		event.BroadcasterUserID,
		pq.StringArray{event.Reward.ID},
	).Find(&alert).Error; err != nil {
		return err
	}

	if alert.ID == "" {
		return nil
	}

	_, err := c.websocketsGrpc.TriggerAlert(
		context.TODO(),
		&websockets.TriggerAlertRequest{
			ChannelId: event.BroadcasterUserID,
			AlertId:   alert.ID,
		},
	)

	return err
}

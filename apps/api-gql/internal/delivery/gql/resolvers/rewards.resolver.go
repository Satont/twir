package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"slices"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

// TwitchRewards is the resolver for the twitchRewards field.
func (r *queryResolver) TwitchRewards(ctx context.Context, channelID *string) ([]gqlmodel.TwitchReward, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	channelIdForRequest := dashboardId
	if channelID != nil {
		channelIdForRequest = *channelID
	}

	rewards, err := r.cachedTwitchClient.GetChannelRewards(ctx, channelIdForRequest)
	if err != nil {
		return nil, err
	}

	var gqlRewards []gqlmodel.TwitchReward
	for _, reward := range rewards {
		imageUrls := append(
			[]string{},
			reward.Image.Url1x,
			reward.Image.Url2x,
			reward.Image.Url4x,
		)

		var usedTimes int64
		if err := r.gorm.
			WithContext(ctx).
			Where("channel_id = ? AND reward_id = ?", channelIdForRequest, reward.ID).
			Model(&model.ChannelRedemption{}).
			Count(&usedTimes).
			Error; err != nil {
			return nil, fmt.Errorf("failed to count channel redemptions: %w", err)
		}

		gqlRewards = append(
			gqlRewards, gqlmodel.TwitchReward{
				ID:              reward.ID,
				Title:           reward.Title,
				Cost:            reward.Cost,
				ImageUrls:       imageUrls,
				BackgroundColor: reward.BackgroundColor,
				Enabled:         false,
				UsedTimes:       int(usedTimes),
			},
		)
	}

	slices.SortFunc(
		gqlRewards,
		func(a, b gqlmodel.TwitchReward) int {
			return b.Cost - a.Cost
		},
	)

	return gqlRewards, nil
}
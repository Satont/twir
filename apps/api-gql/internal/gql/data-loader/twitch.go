package data_loader

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/libs/cache/twitch"
)

func (c *DataLoader) getHelixUsersByIds(ctx context.Context, ids []string) (
	[]*gqlmodel.TwirUserTwitchInfo,
	[]error,
) {
	users, err := c.cachedTwitchClient.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	mappedUsers := make([]*gqlmodel.TwirUserTwitchInfo, 0, len(users))

	for _, id := range ids {
		user, ok := lo.Find(
			users, func(item twitch.TwitchUser) bool {
				return item.ID == id
			},
		)
		if !ok {
			mappedUsers = append(
				mappedUsers, &gqlmodel.TwirUserTwitchInfo{
					ID:          id,
					Login:       "[twir] twitch banned",
					DisplayName: "[Twir] Twitch Banned",
				},
			)
			continue
		}

		mappedUsers = append(
			mappedUsers, &gqlmodel.TwirUserTwitchInfo{
				ID:              user.ID,
				Login:           user.Login,
				DisplayName:     user.DisplayName,
				ProfileImageURL: user.ProfileImageURL,
				Description:     user.Description,
			},
		)
	}

	return mappedUsers, nil
}

func GetHelixUser(ctx context.Context, userID string) (*gqlmodel.TwirUserTwitchInfo, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserLoader.Load(ctx, userID)
}

func GetHelixUsers(ctx context.Context, userIDs []string) ([]*gqlmodel.TwirUserTwitchInfo, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserLoader.LoadAll(ctx, userIDs)
}

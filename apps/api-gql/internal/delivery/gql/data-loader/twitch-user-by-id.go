package data_loader

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/cache/twitch"
)

func (c *DataLoader) getHelixUsersByIds(ctx context.Context, ids []string) (
	[]*gqlmodel.TwirUserTwitchInfo,
	[]error,
) {
	nonEmptyIds := lo.Filter(
		ids, func(id string, _ int) bool {
			return id != ""
		},
	)

	users, err := c.cachedTwitchClient.GetUsersByIds(ctx, nonEmptyIds)
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
					Login:       "[twir] not found",
					DisplayName: "[Twir] Not Found",
					NotFound:    true,
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
				OfflineImageURL: user.OfflineImageURL,
				Description:     user.Description,
			},
		)
	}

	return mappedUsers, nil
}

func GetHelixUserById(ctx context.Context, userID string) (*gqlmodel.TwirUserTwitchInfo, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserByIdLoader.Load(ctx, userID)
}

func GetHelixUsersByIds(ctx context.Context, userIDs []string) (
	[]*gqlmodel.TwirUserTwitchInfo,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserByIdLoader.LoadAll(ctx, userIDs)
}

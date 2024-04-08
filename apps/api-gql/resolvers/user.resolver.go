package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
)

// AuthedUser is the resolver for the authedUser field.
func (r *queryResolver) AuthedUser(ctx context.Context) (*gqlmodel.User, error) {
	user, err := sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("not authenticated")
	}
	return &gqlmodel.User{
		ID:                user.ID,
		IsBotAdmin:        user.IsBotAdmin,
		APIKey:            user.ApiKey,
		IsBanned:          user.IsBanned,
		HideOnLandingPage: user.HideOnLandingPage,
		Channel: &gqlmodel.UserChannel{
			IsEnabled:      user.Channel.IsEnabled,
			IsBotModerator: user.Channel.IsBotMod,
			BotID:          user.Channel.BotID,
		},
	}, nil
}

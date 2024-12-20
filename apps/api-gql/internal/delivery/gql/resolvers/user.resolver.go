package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"gorm.io/gorm"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *authenticatedUserResolver) TwitchProfile(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.ID)
}

// SelectedDashboardTwitchUser is the resolver for the selectedDashboardTwitchUser field.
func (r *authenticatedUserResolver) SelectedDashboardTwitchUser(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.SelectedDashboardID)
}

// AvailableDashboards is the resolver for the availableDashboards field.
func (r *authenticatedUserResolver) AvailableDashboards(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) ([]gqlmodel.Dashboard, error) {
	return r.getAvailableDashboards(ctx, obj)
}

// TwitchProfile is the resolver for the twitchProfile field.
func (r *dashboardResolver) TwitchProfile(
	ctx context.Context,
	obj *gqlmodel.Dashboard,
) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.ID)
}

// AuthenticatedUserSelectDashboard is the resolver for the authenticatedUserSelectDashboard field.
func (r *mutationResolver) AuthenticatedUserSelectDashboard(
	ctx context.Context,
	dashboardID string,
) (bool, error) {
	if err := r.sessions.SetSessionSelectedDashboard(ctx, dashboardID); err != nil {
		return false, err
	}

	return true, nil
}

// AuthenticatedUserUpdateSettings is the resolver for the authenticatedUserUpdateSettings field.
func (r *mutationResolver) AuthenticatedUserUpdateSettings(
	ctx context.Context,
	opts gqlmodel.UserUpdateSettingsInput,
) (bool, error) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := &model.Users{}
	if err := r.gorm.
		WithContext(ctx).
		Where("id = ?", user.ID).
		First(entity).Error; err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if opts.HideOnLandingPage.IsSet() {
		entity.HideOnLandingPage = *opts.HideOnLandingPage.Value()
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(entity).Error; err != nil {
		return false, fmt.Errorf("failed to save user: %w", err)
	}

	return true, nil
}

// AuthenticatedUserRegenerateAPIKey is the resolver for the authenticatedUserRegenerateApiKey field.
func (r *mutationResolver) AuthenticatedUserRegenerateAPIKey(ctx context.Context) (string, error) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return "", err
	}

	entity := &model.Users{}
	if err := r.gorm.
		WithContext(ctx).
		Where("id = ?", user.ID).
		First(entity).Error; err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	entity.ApiKey = uuid.NewString()

	if err := r.gorm.
		WithContext(ctx).
		Save(entity).Error; err != nil {
		return "", fmt.Errorf("failed to save user: %w", err)
	}

	return entity.ApiKey, nil
}

// AuthenticatedUserUpdatePublicPage is the resolver for the authenticatedUserUpdatePublicPage field.
func (r *mutationResolver) AuthenticatedUserUpdatePublicPage(
	ctx context.Context,
	opts gqlmodel.UserUpdatePublicSettingsInput,
) (bool, error) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	currentSettings := &model.ChannelPublicSettings{}
	if err := r.gorm.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			user.ID,
		).
		Preload("SocialLinks").
		// init default settings
		FirstOrInit(
			currentSettings,
			&model.ChannelPublicSettings{
				ChannelID: user.ID,
			},
		).
		Error; err != nil {
		return false, fmt.Errorf("failed to get public settings: %w", err)
	}

	txErr := r.gorm.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if opts.Description.IsSet() {
				currentSettings.Description = null.StringFromPtr(opts.Description.Value())
			}

			if opts.SocialLinks.IsSet() {
				if err := tx.
					Where("settings_id = ?", currentSettings.ID).
					Delete(&model.ChannelPublicSettingsSocialLink{}).
					Error; err != nil {
					return err
				}

				links := make([]model.ChannelPublicSettingsSocialLink, 0, len(opts.SocialLinks.Value()))
				for _, link := range opts.SocialLinks.Value() {
					links = append(
						links,
						model.ChannelPublicSettingsSocialLink{
							ID:         uuid.New(),
							SettingsID: currentSettings.ID,
							Title:      link.Title,
							Href:       link.Href,
						},
					)
				}

				currentSettings.SocialLinks = links
			}

			return tx.Save(currentSettings).Error
		},
	)

	if txErr != nil {
		return false, fmt.Errorf("failed to update public settings: %w", txErr)
	}

	return true, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	if err := r.sessions.SessionLogout(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// AuthenticatedUser is the resolver for the authenticatedUser field.
func (r *queryResolver) AuthenticatedUser(ctx context.Context) (
	*gqlmodel.AuthenticatedUser,
	error,
) {
	sessionUser, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("not authenticated: %w", err)
	}

	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user := model.Users{}
	if err := r.gorm.
		WithContext(ctx).
		Where("id = ?", sessionUser.ID).
		Preload("Channel").
		First(&user).Error; err != nil {
	}

	authedUser := &gqlmodel.AuthenticatedUser{
		ID:                  user.ID,
		IsBotAdmin:          user.IsBotAdmin,
		IsBanned:            user.IsBanned,
		HideOnLandingPage:   user.HideOnLandingPage,
		APIKey:              user.ApiKey,
		SelectedDashboardID: dashboardId,
	}

	if user.Channel != nil {
		authedUser.IsEnabled = &user.Channel.IsEnabled
		authedUser.IsBotModerator = &user.Channel.IsBotMod
		authedUser.BotID = &user.Channel.BotID
	}

	return authedUser, nil
}

// UserPublicSettings is the resolver for the userPublicSettings field.
func (r *queryResolver) UserPublicSettings(
	ctx context.Context,
	userID *string,
) (*gqlmodel.PublicSettings, error) {
	dashboardId, _ := r.sessions.GetSelectedDashboard(ctx)

	var idForFetch string
	if userID != nil {
		idForFetch = *userID
	} else if dashboardId != "" {
		idForFetch = dashboardId
	}

	if idForFetch == "" {
		return nil, fmt.Errorf("id for fetch not setted in request or session")
	}

	entity := &model.ChannelPublicSettings{}
	if err := r.gorm.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			idForFetch,
		).
		Preload("SocialLinks").
		Find(entity).Error; err != nil {
		return nil, err
	}

	settings := &gqlmodel.PublicSettings{
		Description: entity.Description.Ptr(),
		SocialLinks: lo.Map(
			entity.SocialLinks,
			func(item model.ChannelPublicSettingsSocialLink, _ int) gqlmodel.SocialLink {
				return gqlmodel.SocialLink{
					Title: item.Title,
					Href:  item.Href,
				}
			},
		),
	}

	return settings, nil
}

// AuthLink is the resolver for the authLink field.
func (r *queryResolver) AuthLink(ctx context.Context, redirectTo string) (string, error) {
	if redirectTo == "" {
		return "", fmt.Errorf("incorrect auth link %s", redirectTo)
	}

	twitchClient, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:    r.config.TwitchClientId,
			RedirectURI: r.config.TwitchCallbackUrl,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create twitch client: %w", err)
	}

	state := base64.StdEncoding.EncodeToString([]byte(redirectTo))

	url := twitchClient.GetAuthorizationURL(
		&helix.AuthorizationURLParams{
			ResponseType: "code",
			Scopes:       twitchScopes,
			State:        state,
			ForceVerify:  false,
		},
	)

	return url, nil
}

// AuthenticatedUser returns graph.AuthenticatedUserResolver implementation.
func (r *Resolver) AuthenticatedUser() graph.AuthenticatedUserResolver {
	return &authenticatedUserResolver{r}
}

// Dashboard returns graph.DashboardResolver implementation.
func (r *Resolver) Dashboard() graph.DashboardResolver { return &dashboardResolver{r} }

type authenticatedUserResolver struct{ *Resolver }
type dashboardResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var twitchScopes = []string{
	"moderation:read",
	"channel:manage:broadcast",
	"channel:read:redemptions",
	"channel:manage:redemptions",
	"moderator:read:chatters",
	"moderator:manage:shoutouts",
	"moderator:manage:banned_users",
	"channel:read:vips",
	"channel:manage:vips",
	"channel:manage:moderators",
	"moderator:read:followers",
	"moderator:manage:chat_settings",
	"channel:read:polls",
	"channel:manage:polls",
	"channel:read:predictions",
	"channel:manage:predictions",
	"channel:read:subscriptions",
	"channel:moderate",
	"user:read:follows",
	"channel:bot",
	"channel:manage:raids",
}

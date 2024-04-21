package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"gorm.io/gorm"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *authenticatedUserResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.AuthenticatedUser) (*gqlmodel.TwirUserTwitchInfo, error) {
	user, err := data_loader.GetHelixUser(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &gqlmodel.TwirUserTwitchInfo{
		Login:           user.Login,
		DisplayName:     user.DisplayName,
		ProfileImageURL: user.ProfileImageURL,
		Description:     user.Description,
	}, nil
}

// TwitchProfile is the resolver for the twitchProfile field.
func (r *dashboardResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.Dashboard) (*gqlmodel.TwirUserTwitchInfo, error) {
	user, err := data_loader.GetHelixUser(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &gqlmodel.TwirUserTwitchInfo{
		Login:           user.Login,
		DisplayName:     user.DisplayName,
		ProfileImageURL: user.ProfileImageURL,
		Description:     user.Description,
	}, nil
}

// AuthenticatedUserSelectDashboard is the resolver for the authenticatedUserSelectDashboard field.
func (r *mutationResolver) AuthenticatedUserSelectDashboard(ctx context.Context, dashboardID string) (bool, error) {
	if err := r.sessions.SetSelectedDashboard(ctx, dashboardID); err != nil {
		return false, err
	}

	return true, nil
}

// AuthenticatedUserUpdateSettings is the resolver for the authenticatedUserUpdateSettings field.
func (r *mutationResolver) AuthenticatedUserUpdateSettings(ctx context.Context, opts gqlmodel.UserUpdateSettingsInput) (bool, error) {
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
func (r *mutationResolver) AuthenticatedUserUpdatePublicPage(ctx context.Context, opts gqlmodel.UserUpdatePublicSettingsInput) (bool, error) {
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

// AuthenticatedUser is the resolver for the authenticatedUser field.
func (r *queryResolver) AuthenticatedUser(ctx context.Context) (*gqlmodel.AuthenticatedUser, error) {
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
		TwitchProfile:       &gqlmodel.TwirUserTwitchInfo{},
		APIKey:              user.ApiKey,
		SelectedDashboardID: dashboardId,
		AvailableDashboards: []gqlmodel.Dashboard{},
	}

	if user.Channel != nil {
		authedUser.IsEnabled = &user.Channel.IsEnabled
		authedUser.IsBotModerator = &user.Channel.IsBotMod
		authedUser.BotID = &user.Channel.BotID
	}

	var dashboardsEntities []gqlmodel.Dashboard
	if authedUser.IsBotAdmin {
		var channels []model.Channels
		if err := r.gorm.WithContext(ctx).Find(&channels).Error; err != nil {
			return nil, err
		}

		for _, channel := range channels {
			dashboardsEntities = append(
				dashboardsEntities,
				gqlmodel.Dashboard{
					ID: channel.ID,
					Flags: []gqlmodel.ChannelRolePermissionEnum{
						gqlmodel.ChannelRolePermissionEnumCanAccessDashboard,
					},
				},
			)
		}
	} else {
		var roles []model.ChannelRoleUser
		if err := r.gorm.
			WithContext(ctx).
			Where(
				`"userId" = ?`,
				user.ID,
			).
			Preload("Role").
			Preload("Role.Channel").
			Find(&roles).
			Error; err != nil {
			return nil, err
		}
		for _, role := range roles {
			if role.Role == nil || role.Role.Channel == nil || len(role.Role.Permissions) == 0 {
				continue
			}

			var flags []gqlmodel.ChannelRolePermissionEnum
			for _, flag := range role.Role.Permissions {
				flags = append(flags, gqlmodel.ChannelRolePermissionEnum(flag))
			}

			dashboardsEntities = append(
				dashboardsEntities,
				gqlmodel.Dashboard{
					ID:    role.Role.Channel.ID,
					Flags: flags,
				},
			)
		}
	}

	var usersStats []model.UsersStats
	if err := r.gorm.
		WithContext(ctx).
		Where(`"userId" = ?`, user.ID).
		Find(&usersStats).Error; err != nil {
		return nil, err
	}

	for _, stat := range usersStats {
		var channelRoles []model.ChannelRole
		if err := r.gorm.WithContext(ctx).Where(`"channelId" = ?`, stat.ChannelID).Find(&channelRoles).
			Error; err != nil {
			return nil, err
		}

		var role model.ChannelRole

		if stat.IsMod {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeModerator
				},
			)
		} else if stat.IsVip {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeVip
				},
			)
		} else if stat.IsSubscriber {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeSubscriber
				},
			)
		}

		var flags []gqlmodel.ChannelRolePermissionEnum
		for _, flag := range role.Permissions {
			flags = append(flags, gqlmodel.ChannelRolePermissionEnum(flag))
		}

		if role.ID != "" {
			dashboardsEntities = append(
				dashboardsEntities,
				gqlmodel.Dashboard{
					ID:    role.ChannelID,
					Flags: flags,
				},
			)
		}
	}

	authedUser.AvailableDashboards = dashboardsEntities

	entity := &model.ChannelPublicSettings{}
	if err := r.gorm.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			user.ID,
		).
		Preload("SocialLinks").
		Find(entity).Error; err != nil {
		return nil, err
	}

	authedUser.PublicSettings = &gqlmodel.PublicSettings{
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

	return authedUser, nil
}

// AuthenticatedUser returns graph.AuthenticatedUserResolver implementation.
func (r *Resolver) AuthenticatedUser() graph.AuthenticatedUserResolver {
	return &authenticatedUserResolver{r}
}

// Dashboard returns graph.DashboardResolver implementation.
func (r *Resolver) Dashboard() graph.DashboardResolver { return &dashboardResolver{r} }

type authenticatedUserResolver struct{ *Resolver }
type dashboardResolver struct{ *Resolver }

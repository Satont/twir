package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/utils"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/mappers"
	"gorm.io/gorm"
)

// RolesCreate is the resolver for the rolesCreate field.
func (r *mutationResolver) RolesCreate(ctx context.Context, opts gqlmodel.RolesCreateOrUpdateOpts) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	permissions := make([]string, 0, len(opts.Permissions))
	for _, permission := range opts.Permissions {
		permissions = append(permissions, permission.String())
	}

	users := make([]*model.ChannelRoleUser, 0, len(opts.Users))
	for _, userId := range opts.Users {
		users = append(
			users,
			&model.ChannelRoleUser{
				ID:     uuid.New().String(),
				UserID: userId,
			},
		)
	}

	entity := &model.ChannelRole{
		ID:                        uuid.NewString(),
		ChannelID:                 dashboardId,
		Name:                      opts.Name,
		Type:                      model.ChannelRoleEnum(gqlmodel.RoleTypeEnumCustom.String()),
		Permissions:               permissions,
		RequiredWatchTime:         int64(opts.Settings.RequiredWatchTime),
		RequiredMessages:          int32(opts.Settings.RequiredMessages),
		RequiredUsedChannelPoints: int64(opts.Settings.RequiredUserChannelPoints),
		Users:                     users,
	}

	if err := r.gorm.WithContext(ctx).Create(entity).Error; err != nil {
		return false, err
	}

	r.logger.Audit(
		"Role create",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationCreate,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// RolesUpdate is the resolver for the rolesUpdate field.
func (r *mutationResolver) RolesUpdate(ctx context.Context, id string, opts gqlmodel.RolesCreateOrUpdateOpts) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	var rolesCount int64
	if err := r.gorm.
		WithContext(ctx).
		Model(&model.ChannelRole{}).
		Where(`"channelId" = ?`, dashboardId).
		Count(&rolesCount).
		Error; err != nil {
		return false, fmt.Errorf("failed to count roles: %w", err)
	}

	if rolesCount >= 20 {
		return false, fmt.Errorf("maximum number of roles reached")
	}

	entity := &model.ChannelRole{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"id" = ? AND "channelId" = ?`, id, dashboardId).
		First(entity).
		Error; err != nil {
		return false, fmt.Errorf("failed to find role: %w", err)
	}

	var entityCopy model.ChannelRole
	if err := utils.DeepCopy(entity, &entityCopy); err != nil {
		return false, err
	}

	entity.Name = opts.Name
	entity.RequiredWatchTime = int64(opts.Settings.RequiredWatchTime)
	entity.RequiredMessages = int32(opts.Settings.RequiredMessages)
	entity.RequiredUsedChannelPoints = int64(opts.Settings.RequiredUserChannelPoints)

	permissions := make([]string, 0, len(opts.Permissions))
	for _, permission := range opts.Permissions {
		permissions = append(permissions, permission.String())
	}
	entity.Permissions = permissions

	users := make([]*model.ChannelRoleUser, 0, len(opts.Users))
	for _, userId := range opts.Users {
		users = append(
			users,
			&model.ChannelRoleUser{
				ID:     uuid.New().String(),
				UserID: userId,
				RoleID: entity.ID,
			},
		)
	}
	entity.Users = users

	txErr := r.gorm.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Where(
				`"roleId" = ?`,
				entity.ID,
			).Delete(&model.ChannelRoleUser{}).Error; err != nil {
				return err
			}

			err := tx.WithContext(ctx).Save(entity).Error
			return err
		},
	)
	if txErr != nil {
		return false, fmt.Errorf("failed to update role: %w", txErr)
	}

	r.logger.Audit(
		"Role update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationUpdate,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// RolesRemove is the resolver for the rolesRemove field.
func (r *mutationResolver) RolesRemove(ctx context.Context, id string) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := &model.ChannelRole{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"id" = ? AND "channelId" = ?`, id, dashboardId).
		First(entity).
		Error; err != nil {
		return false, fmt.Errorf("failed to find role: %w", err)
	}

	if entity.Type.String() != model.ChannelRoleTypeCustom.String() {
		return false, fmt.Errorf("cannot remove default roles")
	}

	if err := r.gorm.
		WithContext(ctx).
		Delete(entity).
		Error; err != nil {
		return false, err
	}

	r.logger.Audit(
		"Role remove",
		audit.Fields{
			OldValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationDelete,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context) ([]gqlmodel.Role, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelRole
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Group(`"id"`).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}

	res := make([]gqlmodel.Role, 0, len(entities))
	for _, entity := range entities {
		permissions := make([]gqlmodel.ChannelRolePermissionEnum, 0, len(entity.Permissions))
		for _, permission := range entity.Permissions {
			permissions = append(permissions, gqlmodel.ChannelRolePermissionEnum(permission))
		}

		res = append(
			res,
			gqlmodel.Role{
				ID:          entity.ID,
				ChannelID:   entity.ChannelID,
				Name:        entity.Name,
				Type:        gqlmodel.RoleTypeEnum(entity.Type.String()),
				Permissions: permissions,
				Settings: &gqlmodel.RoleSettings{
					RequiredWatchTime:         int(entity.RequiredWatchTime),
					RequiredMessages:          int(entity.RequiredMessages),
					RequiredUserChannelPoints: int(entity.RequiredUsedChannelPoints),
				},
			},
		)
	}

	slices.SortFunc(
		res, func(a, b gqlmodel.Role) int {
			typeIdx := lo.IndexOf(gqlmodel.AllRoleTypeEnum, a.Type)

			return typeIdx - lo.IndexOf(gqlmodel.AllRoleTypeEnum, b.Type)
		},
	)

	return res, nil
}

// Users is the resolver for the users field.
func (r *roleResolver) Users(ctx context.Context, obj *gqlmodel.Role) ([]gqlmodel.TwirUserTwitchInfo, error) {
	var users []model.ChannelRoleUser
	if err := r.gorm.
		WithContext(ctx).
		Where(`"roleId" = ?`, obj.ID).
		Find(&users).
		Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	ids := make([]string, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.UserID)
	}

	profiles, err := data_loader.GetHelixUsersByIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch twitch profiles: %w", err)
	}

	res := make([]gqlmodel.TwirUserTwitchInfo, 0, len(profiles))
	for _, profile := range profiles {
		res = append(res, *profile)
	}

	return res, nil
}

// Role returns graph.RoleResolver implementation.
func (r *Resolver) Role() graph.RoleResolver { return &roleResolver{r} }

type roleResolver struct{ *Resolver }

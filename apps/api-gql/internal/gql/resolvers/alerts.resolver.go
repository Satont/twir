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
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// ChannelAlertsCreate is the resolver for the channelAlertsCreate field.
func (r *mutationResolver) ChannelAlertsCreate(
	ctx context.Context,
	input gqlmodel.ChannelAlertCreateInput,
) (*gqlmodel.ChannelAlert, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	volume := 100
	if input.AudioVolume.IsSet() {
		volume = *input.AudioVolume.Value()
	}

	entity := model.ChannelAlert{
		ID:           uuid.NewString(),
		ChannelID:    dashboardId,
		Name:         input.Name,
		AudioID:      null.StringFromPtr(input.AudioID.Value()),
		AudioVolume:  volume,
		CommandIDS:   input.CommandIds.Value(),
		RewardIDS:    input.RewardIds.Value(),
		GreetingsIDS: input.GreetingsIds.Value(),
		KeywordsIDS:  input.KeywordsIds.Value(),
	}

	if err := r.gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	r.logger.Audit(
		"Channel alert create",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        "channel_alerts",
			OperationType: audit.OperationCreate,
			ObjectID:      &entity.ID,
		},
	)

	return &gqlmodel.ChannelAlert{
		ID:           entity.ID,
		Name:         entity.Name,
		AudioID:      entity.AudioID.Ptr(),
		AudioVolume:  &entity.AudioVolume,
		CommandIds:   entity.CommandIDS,
		RewardIds:    entity.RewardIDS,
		GreetingsIds: entity.GreetingsIDS,
		KeywordsIds:  entity.KeywordsIDS,
	}, nil
}

// ChannelAlertsUpdate is the resolver for the channelAlertsUpdate field.
func (r *mutationResolver) ChannelAlertsUpdate(
	ctx context.Context,
	id string,
	input gqlmodel.ChannelAlertUpdateInput,
) (*gqlmodel.ChannelAlert, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelAlert{}
	if err := r.gorm.
		WithContext(ctx).
		Where(
			`id = ? and "channel_id" = ?`,
			id,
			dashboardId,
		).First(&entity).Error; err != nil {
		return nil, fmt.Errorf("channel alert not found: %w", err)
	}

	var entityCopy model.ChannelAlert
	if err := utils.DeepCopy(&entity, &entityCopy); err != nil {
		return nil, err
	}

	if input.AudioVolume.IsSet() {
		entity.AudioVolume = *input.AudioVolume.Value()
	}

	if input.Name.IsSet() {
		entity.Name = *input.Name.Value()
	}

	if input.AudioID.IsSet() {
		entity.AudioID = null.StringFromPtr(input.AudioID.Value())
	}

	if input.CommandIds.IsSet() {
		entity.CommandIDS = input.CommandIds.Value()
	}

	if input.RewardIds.IsSet() {
		entity.RewardIDS = input.RewardIds.Value()
	}

	if input.GreetingsIds.IsSet() {
		entity.GreetingsIDS = input.GreetingsIds.Value()
	}

	if input.KeywordsIds.IsSet() {
		entity.KeywordsIDS = input.KeywordsIds.Value()
	}

	if err := r.gorm.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	r.logger.Audit(
		"Channel alert update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        "channel_alerts",
			OperationType: audit.OperationUpdate,
			ObjectID:      &entity.ID,
		},
	)

	return &gqlmodel.ChannelAlert{
		ID:           entity.ID,
		Name:         entity.Name,
		AudioID:      entity.AudioID.Ptr(),
		AudioVolume:  &entity.AudioVolume,
		CommandIds:   entity.CommandIDS,
		RewardIds:    entity.RewardIDS,
		GreetingsIds: entity.GreetingsIDS,
		KeywordsIds:  entity.KeywordsIDS,
	}, nil
}

// ChannelAlertsDelete is the resolver for the channelAlertsDelete field.
func (r *mutationResolver) ChannelAlertsDelete(ctx context.Context, id string) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelAlert{}
	if err := r.gorm.
		WithContext(ctx).
		Where(
			`id = ? and "channel_id" = ?`, id,
			dashboardId,
		).First(&entity).Error; err != nil {
		return false, err
	}

	if err := r.gorm.WithContext(ctx).Delete(&entity).Error; err != nil {
		return false, err
	}

	r.logger.Audit(
		"Channel alert delete",
		audit.Fields{
			OldValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        "channel_alerts",
			OperationType: audit.OperationDelete,
			ObjectID:      &entity.ID,
		},
	)

	return true, nil
}

// ChannelAlerts is the resolver for the channelAlerts  field.
func (r *queryResolver) ChannelAlerts(ctx context.Context) ([]gqlmodel.ChannelAlert, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelAlert
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		Order("name desc").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	var result []gqlmodel.ChannelAlert
	for _, entity := range entities {
		result = append(
			result,
			gqlmodel.ChannelAlert{
				ID:           entity.ID,
				Name:         entity.Name,
				AudioID:      entity.AudioID.Ptr(),
				AudioVolume:  &entity.AudioVolume,
				CommandIds:   entity.CommandIDS,
				RewardIds:    entity.RewardIDS,
				GreetingsIds: entity.GreetingsIDS,
				KeywordsIds:  entity.KeywordsIDS,
			},
		)
	}

	return result, nil
}

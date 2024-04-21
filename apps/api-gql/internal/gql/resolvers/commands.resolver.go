package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"gorm.io/gorm"
)

// Responses is the resolver for the responses field.
func (r *commandResolver) Responses(
	ctx context.Context,
	obj *gqlmodel.Command,
) ([]gqlmodel.CommandResponse, error) {
	if obj.Default {
		return []gqlmodel.CommandResponse{}, nil
	}

	var responses []model.ChannelsCommandsResponses
	if err := r.gorm.
		WithContext(ctx).
		Where(`"commandId" = ?`, obj.ID).
		Find(&responses).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch command %s responses: %w", obj.ID, err)
	}

	convertedResponses := make([]gqlmodel.CommandResponse, 0, len(responses))
	for _, response := range responses {
		convertedResponses = append(
			convertedResponses,
			gqlmodel.CommandResponse{
				ID:        response.ID,
				CommandID: response.CommandID,
				Text:      response.Text.String,
				Order:     response.Order,
			},
		)
	}

	return convertedResponses, nil
}

// CommandsCreate is the resolver for the commandsCreate field.
func (r *mutationResolver) CommandsCreate(
	ctx context.Context,
	opts gqlmodel.CommandsCreateOpts,
) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	command := &model.ChannelsCommands{
		ID:           uuid.New().String(),
		Name:         strings.ToLower(opts.Name),
		Cooldown:     null.IntFrom(int64(opts.Cooldown)),
		CooldownType: opts.CooldownType,
		Enabled:      opts.Enabled,
		Aliases: lo.Map(
			lo.IfF(
				opts.Aliases == nil, func() []string {
					return []string{}
				},
			).Else(opts.Aliases),
			func(alias string, _ int) string {
				return strings.TrimSuffix(strings.ToLower(alias), "!")
			},
		),
		Description:               null.StringFrom(opts.Description),
		Visible:                   opts.Visible,
		ChannelID:                 dashboardId,
		Default:                   false,
		DefaultName:               null.String{},
		Module:                    "CUSTOM",
		IsReply:                   opts.IsReply,
		KeepResponsesOrder:        opts.KeepResponsesOrder,
		DeniedUsersIDS:            opts.DeniedUsersIds,
		AllowedUsersIDS:           opts.AllowedUsersIds,
		RolesIDS:                  opts.RolesIds,
		OnlineOnly:                opts.OnlineOnly,
		RequiredWatchTime:         opts.RequiredWatchTime,
		RequiredMessages:          opts.RequiredMessages,
		RequiredUsedChannelPoints: opts.RequiredUsedChannelPoints,
		Responses: make(
			[]*model.ChannelsCommandsResponses,
			0,
			len(opts.Responses),
		),
		GroupID:           null.StringFromPtr(opts.GroupID.Value()),
		CooldownRolesIDs:  opts.CooldownRolesIds,
		EnabledCategories: opts.EnabledCategories,
	}

	for _, res := range opts.Responses {
		if res.Text == "" {
			continue
		}

		command.Responses = append(
			command.Responses, &model.ChannelsCommandsResponses{
				ID:    uuid.New().String(),
				Text:  null.StringFrom(res.Text),
				Order: res.Order,
			},
		)
	}

	err = r.gorm.WithContext(ctx).Create(command).Error
	if err != nil {
		return false, fmt.Errorf("failed to create command: %w", err)
	}

	return true, nil
}

// CommandsUpdate is the resolver for the commandsUpdate field.
func (r *mutationResolver) CommandsUpdate(
	ctx context.Context,
	id string,
	opts gqlmodel.CommandsUpdateOpts,
) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	cmd := &model.ChannelsCommands{}
	if err := r.gorm.
		WithContext(ctx).
		Where(
			`"id" = ? AND "channelId" = ?`, id, dashboardId,
		).
		First(cmd).
		Error; err != nil {
		return false, fmt.Errorf("command not found: %w", err)
	}

	if opts.Name.IsSet() {
		cmd.Name = strings.ToLower(*opts.Name.Value())
	}

	if opts.Cooldown.IsSet() {
		cmd.Cooldown = null.IntFrom(int64(*opts.Cooldown.Value()))
	}

	if opts.CooldownType.IsSet() {
		cmd.CooldownType = *opts.CooldownType.Value()
	}

	if opts.Enabled.IsSet() {
		cmd.Enabled = *opts.Enabled.Value()
	}

	if opts.Aliases.IsSet() {
		cmd.Aliases = lo.Map(
			opts.Aliases.Value(),
			func(alias string, _ int) string {
				return strings.TrimSuffix(strings.ToLower(alias), "!")
			},
		)
	}

	if opts.Description.IsSet() {
		cmd.Description = null.StringFromPtr(opts.Description.Value())
	}

	if opts.Visible.IsSet() {
		cmd.Visible = *opts.Visible.Value()
	}

	if opts.IsReply.IsSet() {
		cmd.IsReply = *opts.IsReply.Value()
	}

	if opts.KeepResponsesOrder.IsSet() {
		cmd.KeepResponsesOrder = *opts.KeepResponsesOrder.Value()
	}

	if opts.AllowedUsersIds.IsSet() {
		cmd.AllowedUsersIDS = opts.AllowedUsersIds.Value()
	}

	if opts.DeniedUsersIds.IsSet() {
		cmd.DeniedUsersIDS = opts.DeniedUsersIds.Value()
	}

	if opts.RolesIds.IsSet() {
		cmd.RolesIDS = opts.RolesIds.Value()
	}

	if opts.OnlineOnly.IsSet() {
		cmd.OnlineOnly = *opts.OnlineOnly.Value()
	}

	if opts.RequiredWatchTime.IsSet() {
		cmd.RequiredWatchTime = *opts.RequiredWatchTime.Value()
	}

	if opts.RequiredMessages.IsSet() {
		cmd.RequiredMessages = *opts.RequiredMessages.Value()
	}

	if opts.RequiredUsedChannelPoints.IsSet() {
		cmd.RequiredUsedChannelPoints = *opts.RequiredUsedChannelPoints.Value()
	}

	if opts.GroupID.IsSet() {
		cmd.GroupID = null.StringFromPtr(opts.GroupID.Value())
	}

	if opts.CooldownRolesIds.IsSet() {
		cmd.CooldownRolesIDs = opts.CooldownRolesIds.Value()
	}

	if opts.EnabledCategories.IsSet() {
		cmd.EnabledCategories = opts.EnabledCategories.Value()
	}

	if opts.Responses.IsSet() {
		cmd.Responses = make([]*model.ChannelsCommandsResponses, 0, len(opts.Responses.Value()))
		for _, res := range opts.Responses.Value() {
			if res.Text == "" {
				continue
			}

			response := &model.ChannelsCommandsResponses{
				Text:      null.StringFrom(res.Text),
				Order:     res.Order,
				CommandID: cmd.ID,
			}

			cmd.Responses = append(cmd.Responses, response)
		}
	}

	txErr := r.gorm.
		WithContext(ctx).
		Transaction(
			func(tx *gorm.DB) error {
				if opts.Responses.IsSet() {
					if err = tx.Delete(
						&model.ChannelsCommandsResponses{},
						`"commandId" = ?`,
						cmd.ID,
					).Error; err != nil {
						return err
					}
				}

				return tx.Save(cmd).Error
			},
		)
	if txErr != nil {
		return false, fmt.Errorf("failed to update command: %w", txErr)
	}

	return true, nil
}

// CommandsRemove is the resolver for the commandsRemove field.
func (r *mutationResolver) CommandsRemove(ctx context.Context, id string) (bool, error) {
	cmd := &model.ChannelsCommands{}

	if err := r.gorm.
		WithContext(ctx).
		Where(`"id" = ?`, id).
		First(cmd).
		Error; err != nil {
		return false, err
	}

	if cmd.Default {
		return false, fmt.Errorf("cannot remove default command")
	}

	if err := r.gorm.WithContext(ctx).Delete(&cmd).Error; err != nil {
		return false, err
	}

	return true, nil
}

// Commands is the resolver for the commands field.
func (r *queryResolver) Commands(ctx context.Context) ([]gqlmodel.Command, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelsCommands
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Preload("Group").
		Order("name ASC").
		Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch commands: %w", err)
	}

	convertedCommands := make([]gqlmodel.Command, 0, len(entities))
	for _, entity := range entities {
		cooldown := entity.Cooldown.Int64
		cooldownInt := int(cooldown)

		converted := gqlmodel.Command{
			ID:                        entity.ID,
			Name:                      entity.Name,
			Description:               entity.Description.String,
			Aliases:                   entity.Aliases,
			Cooldown:                  cooldownInt,
			CooldownType:              entity.CooldownType,
			Enabled:                   entity.Enabled,
			Visible:                   entity.Visible,
			Default:                   entity.Default,
			DefaultName:               entity.DefaultName.Ptr(),
			Module:                    entity.Module,
			IsReply:                   entity.IsReply,
			KeepResponsesOrder:        entity.KeepResponsesOrder,
			DeniedUsersIds:            entity.DeniedUsersIDS,
			AllowedUsersIds:           entity.AllowedUsersIDS,
			RolesIds:                  entity.RolesIDS,
			OnlineOnly:                entity.OnlineOnly,
			CooldownRolesIds:          entity.CooldownRolesIDs,
			EnabledCategories:         entity.EnabledCategories,
			RequiredWatchTime:         entity.RequiredWatchTime,
			RequiredMessages:          entity.RequiredMessages,
			RequiredUsedChannelPoints: entity.RequiredUsedChannelPoints,
		}

		if entity.Group != nil {
			converted.Group = &gqlmodel.CommandGroup{
				ID:    entity.Group.ID,
				Name:  entity.Group.Name,
				Color: entity.Group.Color,
			}
		}

		convertedCommands = append(convertedCommands, converted)
	}
	return convertedCommands, nil
}

// Command returns graph.CommandResolver implementation.
func (r *Resolver) Command() graph.CommandResolver { return &commandResolver{r} }

type commandResolver struct{ *Resolver }

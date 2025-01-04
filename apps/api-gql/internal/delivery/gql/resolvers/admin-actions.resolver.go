package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
)

// DropAllAuthSessions is the resolver for the dropAllAuthSessions field.
func (r *mutationResolver) DropAllAuthSessions(ctx context.Context) (bool, error) {
	if err := r.deps.AdminActionsService.DropAllAuthSessions(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// EventsubSubscribe is the resolver for the eventsubSubscribe field.
func (r *mutationResolver) EventsubSubscribe(ctx context.Context, opts gqlmodel.EventsubSubscribeInput) (bool, error) {
	condition := mappers.ConditionTypeGqlToEntity(opts.Condition)
	if condition == "" {
		return false, fmt.Errorf("unknown condition type")
	}

	err := r.deps.AdminActionsService.EventSubSubscribe(
		ctx,
		admin_actions.EventSubSubscribeInput{
			Type:      opts.Type,
			Version:   opts.Version,
			Condition: condition,
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RescheduleTimers is the resolver for the rescheduleTimers field.
func (r *mutationResolver) RescheduleTimers(ctx context.Context) (bool, error) {
	if err := r.deps.AdminActionsService.RescheduleTimers(); err != nil {
		return false, err
	}

	return true, nil
}

// EventsubInitChannels is the resolver for the eventsubInitChannels field.
func (r *mutationResolver) EventsubInitChannels(ctx context.Context) (bool, error) {
	if err := r.deps.AdminActionsService.EventsubReinitChannels(); err != nil {
		return false, err
	}

	return true, nil
}

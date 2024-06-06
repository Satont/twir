package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// ChatOverlayUpdate is the resolver for the chatOverlayUpdate field.
func (r *mutationResolver) ChatOverlayUpdate(
	ctx context.Context,
	opts gqlmodel.ChatOverlayUpdateOpts,
) (bool, error) {
	return r.updateChatOverlay(ctx, opts)
}

// ChatOverlays is the resolver for the chatOverlays field.
func (r *queryResolver) ChatOverlays(ctx context.Context) ([]gqlmodel.ChatOverlay, error) {
	overlays, err := r.chatOverlays(ctx)
	if err != nil {
		return nil, err
	}

	return overlays, nil
}

// ChatOverlaysByID is the resolver for the chatOverlaysById field.
func (r *queryResolver) ChatOverlaysByID(ctx context.Context, id string) (
	*gqlmodel.ChatOverlay,
	error,
) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	return r.getChatOverlaySettings(ctx, id, dashboardId)
}

// ChatOverlaySettings is the resolver for the chatOverlaySettings field.
func (r *subscriptionResolver) ChatOverlaySettings(
	ctx context.Context,
	id string,
	apiKey string,
) (<-chan *gqlmodel.ChatOverlay, error) {
	user := model.Users{}
	if err := r.gorm.Where(`"apiKey" = ?`, apiKey).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	channel := make(chan *gqlmodel.ChatOverlay)

	go func() {
		sub, err := r.wsRouter.Subscribe(
			[]string{
				chatOverlaySubscriptionKeyCreate(id, user.ID),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		initialSettings, err := r.getChatOverlaySettings(ctx, id, user.ID)
		if err == nil {
			channel <- initialSettings
		}

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var settings gqlmodel.ChatOverlay
				if err := json.Unmarshal(data, &settings); err != nil {
					panic(err)
				}

				channel <- &settings
			}
		}
	}()

	return channel, nil
}

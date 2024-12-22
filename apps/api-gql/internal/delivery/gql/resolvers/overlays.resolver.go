package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

// ChatOverlayUpdate is the resolver for the chatOverlayUpdate field.
func (r *mutationResolver) ChatOverlayUpdate(ctx context.Context, id string, opts gqlmodel.ChatOverlayMutateOpts) (bool, error) {
	return r.updateChatOverlay(ctx, id, opts)
}

// ChatOverlayCreate is the resolver for the chatOverlayCreate field.
func (r *mutationResolver) ChatOverlayCreate(ctx context.Context, opts gqlmodel.ChatOverlayMutateOpts) (bool, error) {
	return r.chatOverlayCreate(ctx, opts)
}

// ChatOverlayDelete is the resolver for the chatOverlayDelete field.
func (r *mutationResolver) ChatOverlayDelete(ctx context.Context, id string) (bool, error) {
	return r.chatOverlayDelete(ctx, id)
}

// NowPlayingOverlayUpdate is the resolver for the nowPlayingOverlayUpdate field.
func (r *mutationResolver) NowPlayingOverlayUpdate(ctx context.Context, id string, opts gqlmodel.NowPlayingOverlayMutateOpts) (bool, error) {
	return r.updateNowPlayingOverlay(ctx, id, opts)
}

// NowPlayingOverlayCreate is the resolver for the nowPlayingOverlayCreate field.
func (r *mutationResolver) NowPlayingOverlayCreate(ctx context.Context, opts gqlmodel.NowPlayingOverlayMutateOpts) (bool, error) {
	return r.createNowPlayingOverlay(ctx, opts)
}

// NowPlayingOverlayDelete is the resolver for the nowPlayingOverlayDelete field.
func (r *mutationResolver) NowPlayingOverlayDelete(ctx context.Context, id string) (bool, error) {
	return r.deleteNowPlayingOverlay(ctx, id)
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
func (r *queryResolver) ChatOverlaysByID(ctx context.Context, id string) (*gqlmodel.ChatOverlay, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	return r.getChatOverlaySettings(ctx, id, dashboardId)
}

// NowPlayingOverlays is the resolver for the nowPlayingOverlays field.
func (r *queryResolver) NowPlayingOverlays(ctx context.Context) ([]gqlmodel.NowPlayingOverlay, error) {
	return r.nowPlayingOverlays(ctx)
}

// NowPlayingOverlaysByID is the resolver for the nowPlayingOverlaysById field.
func (r *queryResolver) NowPlayingOverlaysByID(ctx context.Context, id string) (*gqlmodel.NowPlayingOverlay, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	return r.getNowPlayingOverlaySettings(ctx, id, dashboardID)
}

// ChatOverlaySettings is the resolver for the chatOverlaySettings field.
func (r *subscriptionResolver) ChatOverlaySettings(ctx context.Context, id string, apiKey string) (<-chan *gqlmodel.ChatOverlay, error) {
	user := model.Users{}
	if err := r.deps.Gorm.Where(`"apiKey" = ?`, apiKey).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	channel := make(chan *gqlmodel.ChatOverlay)

	go func() {
		sub, err := r.deps.WsRouter.Subscribe(
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

// NowPlayingOverlaySettings is the resolver for the nowPlayingOverlaySettings field.
func (r *subscriptionResolver) NowPlayingOverlaySettings(ctx context.Context, id string, apiKey string) (<-chan *gqlmodel.NowPlayingOverlay, error) {
	return r.nowPlayingOverlaySettingsSubscription(ctx, id, apiKey)
}

// NowPlayingCurrentTrack is the resolver for the nowPlayingCurrentTrack field.
func (r *subscriptionResolver) NowPlayingCurrentTrack(ctx context.Context, apiKey string) (<-chan *gqlmodel.NowPlayingOverlayTrack, error) {
	return r.nowPlayingCurrentTrackSubscription(ctx, apiKey)
}

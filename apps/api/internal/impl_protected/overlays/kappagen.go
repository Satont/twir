package overlays

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/overlays_kappagen"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

const kappagenOverlayType = "kappagen_overlay"

func (c *Overlays) kappagenDbToGrpc(s model.KappagenOverlaySettings) *overlays_kappagen.Settings {
	return &overlays_kappagen.Settings{
		Emotes: &overlays_kappagen.Settings_Emotes{
			Time:  s.Emotes.Time,
			Max:   s.Emotes.Max,
			Queue: s.Emotes.Queue,
		},
		Size: &overlays_kappagen.Settings_Size{
			RatioNormal: s.Size.RatioNormal,
			RatioSmall:  s.Size.RatioSmall,
			Min:         s.Size.Min,
			Max:         s.Size.Max,
		},
		Cube: &overlays_kappagen.Settings_Cube{
			Speed: s.Cube.Speed,
		},
		Animation: &overlays_kappagen.Settings_Animation{
			FadeIn:  s.Animation.FadeIn,
			FadeOut: s.Animation.FadeOut,
			ZoomIn:  s.Animation.ZoomIn,
			ZoomOut: s.Animation.ZoomOut,
		},
		Animations: lo.Map(
			s.Animations, func(
				v model.KappagenOverlaySettingsAnimationSettings,
				i int,
			) *overlays_kappagen.Settings_AnimationSettings {
				return &overlays_kappagen.Settings_AnimationSettings{
					Style: v.Style,
					Prefs: &overlays_kappagen.Settings_AnimationSettings_Prefs{
						Size:    v.Prefs.Size,
						Center:  v.Prefs.Center,
						Speed:   v.Prefs.Speed,
						Faces:   v.Prefs.Faces,
						Message: v.Prefs.Message,
						Time:    v.Prefs.Time,
					},
					Count:   v.Count,
					Enabled: v.Enabled,
				}
			},
		),
		EnableRave: s.EnableRave,
	}
}

func (c *Overlays) kappagenGrpcToDb(s *overlays_kappagen.Settings) model.KappagenOverlaySettings {
	return model.KappagenOverlaySettings{
		Emotes: model.KappagenOverlaySettingsEmotes{
			Time:  s.Emotes.Time,
			Max:   s.Emotes.Max,
			Queue: s.Emotes.Queue,
		},
		Size: model.KappagenOverlaySettingsSize{
			RatioNormal: s.Size.RatioNormal,
			RatioSmall:  s.Size.RatioSmall,
			Min:         s.Size.Min,
			Max:         s.Size.Max,
		},
		Cube: model.KappagenOverlaySettingsCube{
			Speed: s.Cube.Speed,
		},
		Animation: model.KappagenOverlaySettingsAnimation{
			FadeIn:  s.Animation.FadeIn,
			FadeOut: s.Animation.FadeOut,
			ZoomIn:  s.Animation.ZoomIn,
			ZoomOut: s.Animation.ZoomOut,
		},
		Animations: lo.Map(
			s.Animations, func(
				v *overlays_kappagen.Settings_AnimationSettings,
				i int,
			) model.KappagenOverlaySettingsAnimationSettings {
				return model.KappagenOverlaySettingsAnimationSettings{
					Style: v.Style,
					Prefs: model.KappagenOverlaySettingsAnimationSettingsPrefs{
						Size:    v.Prefs.Size,
						Center:  v.Prefs.Center,
						Speed:   v.Prefs.Speed,
						Faces:   v.Prefs.Faces,
						Message: v.Prefs.Message,
						Time:    v.Prefs.Time,
					},
					Count:   v.Count,
					Enabled: v.Enabled,
				}
			},
		),
		EnableRave: s.EnableRave,
	}
}

func (c *Overlays) OverlayKappaGenGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_kappagen.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, kappagenOverlayType).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	parsedSettings := model.KappagenOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	return c.kappagenDbToGrpc(parsedSettings), nil
}

func (c *Overlays) OverlayKappaGenUpdate(
	ctx context.Context,
	req *overlays_kappagen.Settings,
) (*overlays_kappagen.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, kappagenOverlayType).
		Find(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	if entity.ID == "" {
		entity.ID = uuid.NewString()
		entity.ChannelId = dashboardId
		entity.Type = kappagenOverlayType
	}

	parsedSettings := c.kappagenGrpcToDb(req)
	settingsJson, err := json.Marshal(parsedSettings)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal settings: %w", err)
	}
	entity.Settings = settingsJson
	if err := c.Db.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot update settings: %w", err)
	}

	newSettings := model.KappagenOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &newSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	c.Grpc.Websockets.RefreshKappagenOverlaySettings(
		ctx, &websockets.RefreshKappagenOverlaySettingsRequest{
			ChannelId: dashboardId,
		},
	)

	return c.kappagenDbToGrpc(newSettings), nil
}

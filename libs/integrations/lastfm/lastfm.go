package lastfm

import (
	"errors"

	model "github.com/satont/twir/libs/gomodels"
	api "github.com/shkh/lastfm-go/lastfm"
	"gorm.io/gorm"
)

type Opts struct {
	Gorm        *gorm.DB
	Integration *model.ChannelsIntegrations
}

type Lastfm struct {
	gorm        *gorm.DB
	integration *model.ChannelsIntegrations
	api         *api.Api
	user        *api.UserGetInfo
}

func New(opts Opts) (*Lastfm, error) {
	if !opts.Integration.APIKey.Valid ||
		opts.Integration.Integration == nil ||
		!opts.Integration.Integration.ClientSecret.Valid ||
		!opts.Integration.Integration.APIKey.Valid {
		return nil, errors.New("integration params is not valid")
	}

	lfm := &Lastfm{
		gorm:        opts.Gorm,
		integration: opts.Integration,
		api: api.New(
			opts.Integration.Integration.APIKey.String,
			opts.Integration.Integration.ClientSecret.String,
		),
	}

	lfm.api.SetSession(opts.Integration.APIKey.String)

	user, err := lfm.api.User.GetInfo(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	lfm.user = &user

	return lfm, nil
}

type Track struct {
	Title  string
	Artist string
	Image  string
}

func (c *Lastfm) GetTrack() (*Track, error) {
	tracks, err := c.api.User.GetRecentTracks(
		map[string]interface{}{
			"limit": "1",
			"user":  c.user.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(tracks.Tracks) == 0 || tracks.Tracks[0].NowPlaying != "true" {
		return nil, nil
	}

	track := tracks.Tracks[0]
	var cover string
	if len(track.Images) > 0 {
		cover = track.Images[0].Url
	}

	return &Track{
		Title:  track.Name,
		Artist: track.Artist.Name,
		Image:  cover,
	}, nil
}

package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var ErrIntegrationNotConfigured = errors.New("integration not configured")

type LinksResolverOpts struct {
	fx.In

	Gorm   *gorm.DB
	Config config.Config
}

func NewLinksResolver(opts LinksResolverOpts) *LinksResolver {
	return &LinksResolver{
		gorm:   opts.Gorm,
		config: opts.Config,
	}
}

type LinksResolver struct {
	gorm   *gorm.DB
	config config.Config
}

func (c *LinksResolver) GetIntegrationAuthLink(
	ctx context.Context,
	service gqlmodel.IntegrationService,
) (
	*string,
	error,
) {
	entity := model.Integrations{}
	if err := c.gorm.
		WithContext(ctx).
		Where("service = ?", service).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("integration %s not found: %w", service, err)
	}

	switch service {
	case gqlmodel.IntegrationServiceLastfm:
		return c.lastfm(entity)
	case gqlmodel.IntegrationServiceSpotify:
		return c.spotify(entity)
	case gqlmodel.IntegrationServiceDonationalerts:
		return c.donationAlerts(entity)
	case gqlmodel.IntegrationServiceDiscord:
		return c.discord(entity)
	case gqlmodel.IntegrationServiceVk:
		return c.vk(entity)
	case gqlmodel.IntegrationServiceFaceit:
		return c.faceit(entity)
	case gqlmodel.IntegrationServiceValorant:
		return c.valorant(entity)
	case gqlmodel.IntegrationServiceStreamlabs:
		return c.streamlabs(entity)
	}

	return nil, nil
}

func (c *LinksResolver) lastfm(entity model.Integrations) (*string, error) {
	if !entity.APIKey.Valid || !entity.RedirectURL.Valid {
		return nil, ErrIntegrationNotConfigured
	}

	link := fmt.Sprintf(
		"https://www.last.fm/api/auth/?api_key=%s&cb=%s",
		entity.APIKey.String,
		entity.RedirectURL.String,
	)

	return &link, nil
}

func (c *LinksResolver) spotify(entity model.Integrations) (*string, error) {
	if !entity.ClientID.Valid || !entity.ClientSecret.Valid || !entity.RedirectURL.Valid {
		return nil, ErrIntegrationNotConfigured
	}

	u, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", entity.ClientID.String)
	q.Add("scope", "user-read-currently-playing")
	q.Add("redirect_uri", entity.RedirectURL.String)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

func (c *LinksResolver) donationAlerts(entity model.Integrations) (*string, error) {
	if !entity.ClientID.Valid || !entity.ClientSecret.Valid || !entity.RedirectURL.Valid {
		return nil, ErrIntegrationNotConfigured
	}

	u, _ := url.Parse("https://www.donationalerts.com/oauth/authorize")

	q := u.Query()
	q.Add("client_id", entity.ClientID.String)
	q.Add("response_type", "code")
	q.Add("scope", "oauth-user-show oauth-donation-subscribe")
	q.Add("redirect_uri", entity.RedirectURL.String)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

func (c *LinksResolver) discord(entity model.Integrations) (*string, error) {
	u, _ := url.Parse("https://discord.com/oauth2/authorize")

	if c.config.DiscordClientID == "" || c.config.DiscordClientSecret == "" {
		return nil, errors.New("discord not enabled on our side, please be patient")
	}

	redirectUrl := fmt.Sprintf("http://%s/dashboard/integrations/discord", c.config.SiteBaseUrl)

	q := u.Query()
	q.Add("client_id", c.config.DiscordClientID)
	q.Add("response_type", "code")
	q.Add("permissions", "1497333180438")
	q.Add("scope", "bot applications.commands")
	q.Add("redirect_uri", redirectUrl)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

func (c *LinksResolver) vk(entity model.Integrations) (*string, error) {
	if !entity.ClientID.Valid || !entity.ClientSecret.Valid || !entity.RedirectURL.Valid {
		return nil, errors.New("vk not enabled on our side, please be patient")
	}

	u, _ := url.Parse("https://oauth.vk.com/authorize")

	q := u.Query()
	q.Add("client_id", entity.ClientID.String)
	q.Add("display", "page")
	q.Add("response_type", "code")
	q.Add("scope", "status offline")
	q.Add("redirect_uri", entity.RedirectURL.String)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

func (c *LinksResolver) faceit(entity model.Integrations) (*string, error) {
	if !entity.ClientID.Valid || !entity.ClientSecret.Valid || !entity.RedirectURL.Valid {
		return nil, errors.New("faceit not enabled on our side, please be patient")
	}

	u, _ := url.Parse("https://cdn.faceit.com/widgets/sso/index.html")

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", entity.ClientID.String)
	q.Add("redirect_popup", entity.RedirectURL.String)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

func (c *LinksResolver) valorant(entity model.Integrations) (*string, error) {
	if !entity.ClientID.Valid || !entity.ClientSecret.Valid || !entity.RedirectURL.Valid {
		return nil, errors.New("valorant not enabled on our side, please be patient")
	}

	u, _ := url.Parse("https://auth.riotgames.com/authorize")

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", entity.ClientID.String)
	q.Add("scope", strings.Join([]string{"openid", "offline_access", "cpid"}, "+"))
	q.Add("redirect_uri", entity.RedirectURL.String)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

func (c *LinksResolver) streamlabs(entity model.Integrations) (*string, error) {
	if !entity.ClientID.Valid || !entity.ClientSecret.Valid || !entity.RedirectURL.Valid {
		return nil, errors.New("streamlabs not enabled on our side, please be patient")
	}

	u, _ := url.Parse("https://www.streamlabs.com/api/v2.0/authorize")

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", entity.ClientID.String)
	q.Add("scope", "socket.token donations.read")
	q.Add("redirect_uri", entity.RedirectURL.String)
	u.RawQuery = q.Encode()

	link := u.String()

	return &link, nil
}

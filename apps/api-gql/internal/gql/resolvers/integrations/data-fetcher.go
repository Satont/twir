package integrations

import (
	"context"
	"fmt"

	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DataFetcherOpts struct {
	fx.In

	Gorm              *gorm.DB
	Config            config.Config
	AuthLinksResolver *LinksResolver
}

func NewDataFetcher(opts DataFetcherOpts) *DataFetcher {
	return &DataFetcher{
		gorm:          opts.Gorm,
		config:        opts.Config,
		linksResolver: opts.AuthLinksResolver,
	}
}

type DataFetcher struct {
	gorm          *gorm.DB
	config        config.Config
	linksResolver *LinksResolver
}

func (c *DataFetcher) GetIntegrationData(
	ctx context.Context,
	channelId string,
	service gqlmodel.IntegrationService,
) (gqlmodel.IntegrationData, error) {
	integration, err := c.getChannelIntegration(ctx, channelId, service)
	if err != nil {
		return nil, err
	}

	authLink, _ := c.linksResolver.GetIntegrationAuthLink(ctx, service)
	if err != nil {
		return nil, err
	}

	if integration == nil {
		return nil, nil
	}

	switch service {
	case gqlmodel.IntegrationServiceLastfm:
		result := &gqlmodel.IntegrationDataLastfm{
			AuthLink: authLink,
		}

		if integration.Data != nil && integration.Data.UserName != nil && integration.Data.
			Avatar == nil {
			result.Username = integration.Data.UserName
			result.Avatar = integration.Data.Avatar
		}

		return result, nil
	case gqlmodel.IntegrationServiceSpotify:
		result := &gqlmodel.IntegrationDataSpotify{
			AuthLink: authLink,
		}

		if integration.Data != nil && integration.Data.UserName != nil && integration.Data.
			Avatar != nil {
			result.Username = integration.Data.UserName
			result.Avatar = integration.Data.Avatar
		}

		return result, nil
	case gqlmodel.IntegrationServiceDonationalerts:
		result := &gqlmodel.IntegrationDataDonationAlerts{
			AuthLink: authLink,
		}

		if integration.Data != nil && integration.Data.Name != nil && integration.Data.
			Avatar != nil {
			result.Username = integration.Data.Name
			result.Avatar = integration.Data.Avatar
		}

		return result, nil
	case gqlmodel.IntegrationServiceValorant:
		result := &gqlmodel.IntegrationDataValorant{
			AuthLink: authLink,
		}

		if integration.Data != nil && integration.Data.UserName != nil {
			result.Username = integration.Data.UserName
		}

		return result, nil
	case gqlmodel.IntegrationServiceStreamlabs:
		result := &gqlmodel.IntegrationDataStreamLabs{
			AuthLink: authLink,
		}

		if integration.Data != nil && integration.Data.UserName != nil && integration.Data.
			Avatar != nil {
			result.Username = integration.Data.UserName
			result.Avatar = integration.Data.Avatar
		}

		return result, nil
	case gqlmodel.IntegrationServiceVk:
		result := &gqlmodel.IntegrationDataVk{
			AuthLink: authLink,
		}

		if integration.Data != nil && integration.Data.UserName != nil && integration.Data.
			Avatar != nil {
			result.Username = integration.Data.UserName
			result.Avatar = integration.Data.Avatar
		}

		return result, nil
	case gqlmodel.IntegrationServiceDiscord:
		result := &gqlmodel.IntegrationDataDiscord{
			AuthLink: authLink,
		}

		if integration.Data != nil {
		}

		return result, nil
	}

	return nil, nil
}

func (c *DataFetcher) getChannelIntegration(
	ctx context.Context,
	dashboardId string,
	service gqlmodel.IntegrationService,
) (*model.ChannelsIntegrations, error) {
	entity := model.ChannelsIntegrations{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "Integration"."service" = ?`, dashboardId, service).
		Joins("Integration").
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("integration %s not found: %w", service, err)
	}

	return &entity, nil
}

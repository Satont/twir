package integrations

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_lastfm"
	lfm "github.com/shkh/lastfm-go/lastfm"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsLastFMGetAuthLink(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_lastfm.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceLastfm)
	if err != nil {
		return nil, err
	}

	if !integration.APIKey.Valid || !integration.RedirectURL.Valid {
		return nil, fmt.Errorf("lastfm integration not configured")
	}

	link := fmt.Sprintf(
		"https://www.last.fm/api/auth/?api_key=%s&cb=%s",
		integration.APIKey.String,
		integration.RedirectURL.String,
	)

	return &integrations_lastfm.GetAuthLink{
		Link: link,
	}, nil
}

func (c *Integrations) IntegrationsLastFMGetData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_lastfm.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceLastfm, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_lastfm.GetDataResponse{
		UserName: integration.Data.Name,
		Avatar:   integration.Data.Avatar,
	}, nil
}

func (c *Integrations) IntegrationsLastFMPostCode(
	ctx context.Context,
	request *integrations_lastfm.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceLastfm, dashboardId)
	if err != nil {
		return nil, err
	}

	api := lfm.New(
		integration.Integration.APIKey.String,
		integration.Integration.ClientSecret.String,
	)
	err = api.LoginWithToken(request.Code)
	sessionKey := api.GetSessionKey()

	info, err := api.User.GetInfo(make(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	integration.Data = &model.ChannelsIntegrationsData{
		UserName: &info.Name,
		Avatar:   &info.Images[len(info.Images)-1].Url,
	}
	integration.APIKey = null.StringFrom(sessionKey)

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsLastFMLogout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceLastfm, dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data = nil
	integration.APIKey = null.String{}

	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

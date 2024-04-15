package integrations

import (
	"context"

	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type PostCodeHandlerOpts struct {
	fx.In

	Gorm   *gorm.DB
	Config config.Config
}

func NewPostCodeHandler(opts PostCodeHandlerOpts) *PostCodeHandler {
	return &PostCodeHandler{
		gorm:   opts.Gorm,
		config: opts.Config,
	}
}

type PostCodeHandler struct {
	gorm   *gorm.DB
	config config.Config
}

func (h *PostCodeHandler) PostCode(
	ctx context.Context,
	service gqlmodel.IntegrationService,
	channelId string,
	code string,
) error {
	return nil
}

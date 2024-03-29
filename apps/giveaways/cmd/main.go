package main

import (
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/giveaways/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.FxDiOnlyErrors,
		app.App,
	).Run()
}

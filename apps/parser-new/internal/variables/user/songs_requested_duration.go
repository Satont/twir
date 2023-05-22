package user

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
)

type songsRequestedDurationSumResult struct {
	Sum int64 //or int ,or some else
}

var SongsRequestedDuration = &types.Variable{
	Name:         "user.songs.requested.duration",
	Description:  lo.ToPtr("Duration of requested by user songs"),
	CommandsOnly: true,
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		sum := &songsRequestedDurationSumResult{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Table("channels_requested_songs").
			Select("sum(duration) as sum").
			Where(`"channelId" = ? AND "orderedById" = ?`, parseCtx.Channel.ID, parseCtx.Sender.ID).
			Scan(&sum).
			Error

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = "0"
			return result, nil
		}

		f := time.Duration(sum.Sum) * time.Millisecond
		result.Result = fmt.Sprintf("%.1fh", f.Hours())

		return result, nil
	},
}

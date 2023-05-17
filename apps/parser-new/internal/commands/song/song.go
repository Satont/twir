package song

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	currentsong "github.com/satont/tsuwari/apps/parser/internal/variables/song"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var CurrentSong = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "song",
		Description: null.StringFrom(*currentsong.Variable.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "SONGS",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{"currentsong"},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					currentsong.Variable.Name,
				),
			},
		}

		return result
	},
}

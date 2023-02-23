package dota

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"strconv"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	steamid "github.com/leighmacdonald/steamid/v2/steamid"
	"github.com/samber/lo"
)

var AddAccCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "dota addacc",
		Description: lo.ToPtr("Add dota account for watching games"),
		RolesNames:  []model.ChannelRoleEnum{model.ChannelRoleTypeBroadcaster},
		Visible:     false,
		Module:      lo.ToPtr("DOTA"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}
		db := do.MustInvoke[gorm.DB](di.Provider)

		acc, err := strconv.ParseUint(*ctx.Text, 10, 64)
		if err != nil {
			result.Result = append(result.Result, WRONG_ACCOUNT_ID)
			return result
		}

		ok := lo.Try(func() error {
			n := steamid.SID32(acc)
			steamid.SID32ToSID(n)
			return nil
		})

		if !ok {
			result.Result = append(result.Result, WRONG_ACCOUNT_ID)
			return result
		}

		accId := steamid.SID32(acc)

		var count int64 = 0
		err = db.
			Table("channels_dota_accounts").
			Where(`"channelId" = ? AND "id" = ?`, ctx.ChannelId, strconv.Itoa(int(accId))).
			Count(&count).Error

		if err != nil {
			fmt.Println(err)
			result.Result = append(result.Result, "Error happend on our side.")
			return result
		}

		if count != 0 {
			result.Result = append(result.Result, "Account already added.")
			return result
		}

		err = db.
			Create(&model.ChannelsDotaAccounts{
				ID:        strconv.Itoa(int(accId)),
				ChannelID: ctx.ChannelId,
			}).Error

		if err != nil {
			fmt.Println(err)
			result.Result = append(
				result.Result,
				"Something went wrong on out side when inserting account into db.",
			)
			return result
		}

		result.Result = append(result.Result, "Account added.")
		return result
	},
}

package giveaways

import (
	"context"
	"errors"
	"time"

	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/giveaways"
	giveawaysService "github.com/twirapp/twir/libs/grpc/giveaways"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Giveaways struct {
	*impl_deps.Deps
}

func (g *Giveaways) convertEntity(giveaway *model.ChannelGiveaway) *giveaways.Giveaway {
	return &giveaways.Giveaway{
		Id:                        giveaway.ID,
		Description:               giveaway.Description.ValueOrZero(),
		ChannelId:                 giveaway.ChannelID,
		RequiredMinWatchTime:      int32(giveaway.RequiredMinWatchTime),
		RequiredMinFollowTime:     int32(giveaway.RequiredMinFollowTime),
		RequiredMinMessages:       int32(giveaway.RequiredMinMessages),
		Keyword:                   giveaway.Keyword,
		RequiredMinSubscriberTier: int32(giveaway.RequiredMinSubscriberTier),
		RequiredMinSubscriberTime: int32(giveaway.RequiredMinSubscriberTime),
		RolesIds:                  giveaway.RolesIDS,
		FollowersLuck:             int32(giveaway.FollowersLuck),
		SubscribersLuck:           int32(giveaway.SubscribersLuck),
		CreatedAt:                 giveaway.CreatedAt.String(),
		FinishedAt:                nil,
		IsRunning:                 giveaway.IsRunning,
		FollowersAgeLuck:          giveaway.FollowersAgeLuck,
		WinnersCount:              int32(giveaway.WinnersCount),
		IsFinished:                giveaway.IsFinished,
	}
}

func (g *Giveaways) GiveawaysGetParticipants(
	ctx context.Context,
	req *giveaways.GetParticipantsRequest,
) (*giveaways.GetParticipantsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaway model.ChannelGiveaway
	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetGiveawayId()).
		Group(`"id"`).
		First(&dbGiveaway).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var participants []*model.ChannelGiveawayParticipant
	err = g.Db.WithContext(ctx).
		Where(`"giveaway_id" = ? AND "display_name" LIKE ?`, req.GetGiveawayId(), "%"+req.GetQuery()+"%").
		Find(&participants).
		Error
	if err != nil {
		return nil, err
	}

	var count int64
	err = g.Db.WithContext(ctx).
		Where(`"giveaway_id" = ?`, req.GetGiveawayId()).
		Model(&model.ChannelGiveawayParticipant{}).
		Count(&count).
		Error
	if err != nil {
		return nil, err
	}

	var convertedPaticipants []*giveaways.Winner
	for _, participant := range participants {
		convertedPaticipants = append(convertedPaticipants, &giveaways.Winner{
			UserId:      participant.UserID,
			DisplayName: participant.DisplayName,
		})
	}

	return &giveaways.GetParticipantsResponse{
		Winners:    convertedPaticipants,
		TotalCount: count,
	}, nil
}

func (g *Giveaways) GiveawaysCreate(
	ctx context.Context,
	req *giveaways.CreateRequest,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := &model.ChannelGiveaway{
		ChannelID:                 dashboardId,
		Description:               null.NewString(req.GetDescription(), req.GetDescription() != ""),
		Keyword:                   req.GetKeyword(),
		RequiredMinMessages:       int(req.GetRequiredMinMessages()),
		FollowersLuck:             int(req.GetFollowersLuck()),
		SubscribersLuck:           int(req.GetSubscribersLuck()),
		WinnersCount:              int(req.GetWinnersCount()),
		IsRunning:                 false,
		RequiredMinWatchTime:      int(req.GetRequiredMinWatchTime()),
		RequiredMinFollowTime:     int(req.GetRequiredMinFollowTime()),
		RequiredMinSubscriberTier: int(req.GetRequiredMinSubscriberTier()),
		FinishedAt:                null.NewTime(time.Now(), false),
		RequiredMinSubscriberTime: int(req.GetRequiredMinSubscriberTime()),
		RolesIDS:                  req.GetRolesIds(),
		FollowersAgeLuck:          req.GetFollowersAgeLuck(),
		SubscribersTier1Luck:      0,
		SubscribersTier2Luck:      0,
		SubscribersTier3Luck:      0,
		IsFinished:                false,
	}
	err := g.Db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return g.convertEntity(entity), nil
}

func (g *Giveaways) GiveawaysGetAll(
	ctx context.Context,
	req *emptypb.Empty,
) (*giveaways.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaways []*model.ChannelGiveaway
	err := g.Db.WithContext(ctx).Where(`"channel_id" = ?`, dashboardId).Find(&dbGiveaways).Error
	if err != nil {
		return nil, err
	}

	return &giveaways.GetAllResponse{
		Giveaways: lo.Map(
			dbGiveaways,
			func(giveaway *model.ChannelGiveaway, _ int) *giveaways.Giveaway {
				return g.convertEntity(giveaway)
			},
		),
	}, nil
}

func (g *Giveaways) GiveawaysUpdate(
	ctx context.Context,
	req *giveaways.UpdateRequest,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	dbGiveaway := model.ChannelGiveaway{}
	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetId()).
		First(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	dbGiveaway.Description = null.NewString(req.GetDescription(), req.GetDescription() != "")
	dbGiveaway.Keyword = req.GetKeyword()
	dbGiveaway.RequiredMinMessages = int(req.GetRequiredMinMessages())
	dbGiveaway.FollowersLuck = int(req.GetFollowersLuck())
	dbGiveaway.SubscribersLuck = int(req.GetSubscribersLuck())
	dbGiveaway.WinnersCount = int(req.GetWinnersCount())
	dbGiveaway.RequiredMinWatchTime = int(req.GetRequiredMinWatchTime())
	dbGiveaway.RequiredMinFollowTime = int(req.GetRequiredMinFollowTime())
	dbGiveaway.RequiredMinSubscriberTier = int(req.GetRequiredMinSubscriberTier())
	dbGiveaway.RequiredMinSubscriberTime = int(req.GetRequiredMinSubscriberTime())
	dbGiveaway.RolesIDS = req.GetRolesIds()
	dbGiveaway.FollowersAgeLuck = req.GetFollowersAgeLuck()
	dbGiveaway.IsFinished = req.GetIsFinished()
	dbGiveaway.IsRunning = req.GetIsRunning()

	err = g.Db.WithContext(ctx).Model(&dbGiveaway).Updates(
		map[string]interface{}{
			"Description":               dbGiveaway.Description,
			"Keyword":                   dbGiveaway.Keyword,
			"RequiredMinMessages":       dbGiveaway.RequiredMinMessages,
			"FollowersLuck":             dbGiveaway.FollowersLuck,
			"SubscribersLuck":           dbGiveaway.SubscribersLuck,
			"WinnersCount":              dbGiveaway.WinnersCount,
			"RequiredMinWatchTime":      dbGiveaway.RequiredMinWatchTime,
			"RequiredMinFollowTime":     dbGiveaway.RequiredMinFollowTime,
			"RequiredMinSubscriberTier": dbGiveaway.RequiredMinSubscriberTier,
			"RequiredMinSubscriberTime": dbGiveaway.RequiredMinSubscriberTime,
			"FollowersAgeLuck":          dbGiveaway.FollowersAgeLuck,
			"IsFinished":                dbGiveaway.IsFinished,
			"IsRunning":                 dbGiveaway.IsRunning,
		},
	).Error
	if err != nil {
		return nil, err
	}

	return g.convertEntity(&dbGiveaway), nil
}

func (g *Giveaways) GiveawaysDelete(
	ctx context.Context,
	req *giveaways.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetId()).
		Delete(&model.ChannelGiveaway{}).
		Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *Giveaways) GiveawaysGetById(
	ctx context.Context,
	req *giveaways.GetByIdRequest,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaway model.ChannelGiveaway
	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetId()).
		Group(`"id`).
		First(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	return g.convertEntity(&dbGiveaway), nil
}

func (g *Giveaways) GiveawaysChooseWinners(
	ctx context.Context,
	req *giveaways.ChooseWinnersRequest,
) (*giveaways.ChooseWinnersResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaway model.ChannelGiveaway
	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetGiveawayId()).
		Group(`"id`).
		First(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	res, err := g.Grpc.Giveaways.ChooseWinner(ctx, &giveawaysService.ChooseWinnerRequest{
		GiveawayId: req.GetGiveawayId(),
	})
	if err != nil {
		return nil, err
	}

	winners := make([]*giveaways.Winner, len(res.GetWinners()))
	for i, winner := range res.GetWinners() {
		winners[i] = &giveaways.Winner{
			UserId:      winner.GetUserId(),
			DisplayName: winner.GetDisplayName(),
		}
	}

	return &giveaways.ChooseWinnersResponse{
		Winners: winners,
	}, nil
}

func (g *Giveaways) GiveawaysGetWinners(
	ctx context.Context,
	req *giveaways.GetWinnersRequest,
) (*giveaways.GetWinnersResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaway model.ChannelGiveaway
	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetGiveawayId()).
		Group(`"id`).
		First(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	var winners []*model.ChannelGiveawayParticipant
	err = g.Db.WithContext(ctx).
		Where(`"giveaway_id" = ? AND "is_winner" = ?`, req.GetGiveawayId(), true).
		Find(&winners).
		Error
	if err != nil {
		return nil, err
	}

	var convertedWinners []*giveaways.Winner
	for _, winner := range winners {
		convertedWinners = append(convertedWinners, &giveaways.Winner{
			UserId:      winner.UserID,
			DisplayName: winner.DisplayName,
		})
	}

	return &giveaways.GetWinnersResponse{
		Winners: convertedWinners,
	}, nil
}

func (g *Giveaways) GiveawaysClearParticipants(
	ctx context.Context,
	req *giveaways.ClearParticipantsRequest,
) (*giveaways.ClearParticipantsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var giveaway model.ChannelGiveaway
	err := g.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetGiveawayId()).
		First(&giveaway).
		Error
	if err != nil {
		return nil, err
	}

	err = g.Db.WithContext(ctx).
		Where(`"giveaway_id" = ?`, req.GetGiveawayId()).
		Delete(&model.ChannelGiveawayParticipant{}).
		Error
	if err != nil {
		return nil, err
	}

	return &giveaways.ClearParticipantsResponse{
		Winners:    []*giveaways.Winner{},
		TotalCount: 0,
	}, nil
}

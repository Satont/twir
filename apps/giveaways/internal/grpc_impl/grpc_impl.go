package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/gogo/status"
	"github.com/guregu/null"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/giveaways"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Cfg    cfg.Config
	Db     *gorm.DB
	Redis  *redis.Client

	TokensGrpc tokens.TokensClient
}

func New(opts Opts) error {
	impl := &GiveawaysGrpcImplementation{
		db:         opts.Db,
		redis:      opts.Redis,
		logger:     opts.Logger,
		cfg:        opts.Cfg,
		tokensGrpc: opts.TokensGrpc,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.GIVEAWAYS_SERVER_PORT))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
	)

	giveaways.RegisterGiveawaysServer(grpcServer, impl)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(lis)
				opts.Logger.Info(
					"Grpc server started",
					slog.Int("port", constants.GIVEAWAYS_SERVER_PORT),
				)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return nil
}

type GiveawaysGrpcImplementation struct {
	giveaways.UnimplementedGiveawaysServer

	db     *gorm.DB
	redis  *redis.Client
	logger logger.Logger
	cfg    cfg.Config

	tokensGrpc tokens.TokensClient
}

func (c *GiveawaysGrpcImplementation) TryProcessParticipant(
	ctx context.Context,
	req *giveaways.TryProcessParticipantRequest,
) (*emptypb.Empty, error) {
	giveaways := []*model.ChannelGiveaway{}
	err := c.db.WithContext(ctx).
		Where(`"channel_id" = ? AND "is_finished" = ? AND "is_running" = ?`, req.GetChannelId(), false, true).
		Find(&giveaways).
		Error
	if err != nil {
		c.logger.Error(
			"cannot get giveaway",
			slog.Any("err", err),
			slog.String("channelId", req.GetChannelId()),
			slog.String("userId", req.GetUserId()),
		)

		return nil, err
	}

	if len(giveaways) == 0 {
		return &emptypb.Empty{}, nil
	}

	for _, giveaway := range giveaways {
		if !strings.Contains(req.GetMessageText(), giveaway.Keyword) {
			continue
		}

		var dbUser model.Users
		err = c.db.WithContext(ctx).Where(`"id" = ?`, req.GetUserId()).First(&dbUser).Error
		if err != nil {
			c.logger.Error("Cannot get user", slog.Any("err", err))
			return nil, err
		}

		var roles []*model.ChannelRoleUser
		err = c.db.WithContext(ctx).
			Where(`"userId" = ?`, dbUser.ID).
			Preload("Role").
			Find(&roles).
			Error
		if err != nil {
			c.logger.Error("Cannot get user roles", slog.Any("err", err))
			return nil, err
		}
		c.logger.Info("roles", slog.Any("roles", roles))

		var userStats model.UsersStats
		err = c.db.Where(`"userId" = ? AND "channelId" = ?`, dbUser.ID, req.GetChannelId()).
			First(&userStats).
			Error
		if err != nil {
			c.logger.Error("Cannot get user stats", slog.Any("err", err))
			return nil, err
		}

		// TODO: check if user has all roles, and all required stats
		if userStats.Messages < int32(giveaway.RequiredMinMessages) {
			return &emptypb.Empty{}, nil
		}

		if userStats.Watched < int64(giveaway.RequiredMinWatchTime) {
			return &emptypb.Empty{}, nil
		}

		isFollower, followedTime, err := c.isFollower(ctx, req.GetChannelId(), dbUser.ID)
		if err != nil {
			c.logger.Error("Cannot check if user is follower", slog.Any("err", err))
			return nil, err
		}

		newParticipant := model.ChannelGiveawayParticipant{
			GiveawayID:           giveaway.ID,
			UserID:               dbUser.ID,
			DisplayName:          req.GetDisplayName(),
			IsSubscriber:         userStats.IsSubscriber,
			IsModerator:          userStats.IsMod,
			IsVip:                userStats.IsVip,
			IsFollower:           isFollower,
			IsWinner:             false,
			SubscriberTier:       null.IntFrom(0),
			UserFollowSince:      null.NewTime(lo.FromPtr(followedTime), followedTime != nil),
			UserStatsWatchedTime: userStats.Watched,
		}

		err = c.db.WithContext(ctx).Create(&newParticipant).Error
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return &emptypb.Empty{}, nil
				}
			}
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (c *GiveawaysGrpcImplementation) ChooseWinner(
	ctx context.Context,
	req *giveaways.ChooseWinnerRequest,
) (*giveaways.ChooseWinnerResponse, error) {
	giveaway := model.ChannelGiveaway{}

	err := c.db.WithContext(ctx).
		Where(`"id" = ? AND "is_finished" = ?`, req.GetGiveawayId(), false).
		Find(&giveaway).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Cannot find giveaway with this id")
		}

		c.logger.Error(
			"cannot get giveaway",
			slog.Any("err", err),
		)

		return nil, err
	}

	var participants []*model.ChannelGiveawayParticipant
	err = c.db.WithContext(ctx).
		Where(`"giveaway_id" = ?`, req.GetGiveawayId()).
		Find(&participants).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Cannot find giveaway with this id")
		}

		return nil, err
	}

	if len(participants) == 0 {
		return nil, status.Error(codes.Canceled, "No participants")
	}

	if len(participants) <= giveaway.WinnersCount {
		return nil, status.Error(
			codes.OutOfRange,
			"Participants count must be greater than winners count",
		)
	}

	for _, participant := range participants {
		err = c.db.WithContext(ctx).
			Model(participant).
			Update("is_winner", false).
			Error
		if err != nil {
			return nil, err
		}
	}

	processedParticipants := make([]*giveaways.SimplifiedWinner, 0, len(participants))
	for _, participant := range participants {
		countOfTimes := 1
		if participant.IsSubscriber {
			countOfTimes += giveaway.SubscribersLuck
		}
		if participant.IsFollower {
			countOfTimes += giveaway.FollowersLuck
		}

		for i := 0; i < countOfTimes; i++ {
			processedParticipants = append(processedParticipants, &giveaways.SimplifiedWinner{
				UserId:      participant.UserID,
				DisplayName: participant.DisplayName,
			})
		}
	}

	winners := make([]*giveaways.SimplifiedWinner, giveaway.WinnersCount)
	err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(winners); i++ {
			randInd := rand.Intn(len(processedParticipants))
			winners[i] = processedParticipants[randInd]
			processedParticipants[randInd] = processedParticipants[len(processedParticipants)-1]
			processedParticipants = processedParticipants[:len(processedParticipants)-1]

			err = c.db.WithContext(ctx).
				Where(`"giveaway_id" = ? AND "user_id" = ?`, giveaway.ID, winners[i].UserId).
				Model(&model.ChannelGiveawayParticipant{}).
				Update("is_winner", true).Error
			if err != nil {
				return err
			}
		}

		err = c.db.WithContext(ctx).
			Where(`"id" = ?`, giveaway.ID).
			Model(&model.ChannelGiveaway{}).
			Update("is_running", false).
			Error

		return err
	})
	if err != nil {
		return nil, err
	}

	return &giveaways.ChooseWinnerResponse{
		Winners: winners,
	}, nil
}

func (c *GiveawaysGrpcImplementation) subscriberInfo(
	ctx context.Context,
	channelId, userId string,
) (bool, int, error) {
	channel := model.Channels{}
	err := c.db.
		WithContext(ctx).
		Where(`"id" = ?`, channelId).
		First(&channel).
		Error
	if err != nil {
		c.logger.Info("cannot get channel", slog.String("channelId", channelId))
		return false, 0, err
	}

	twitchClient, err := twitch.NewBotClientWithContext(
		ctx,
		channel.BotID,
		c.cfg,
		c.tokensGrpc,
	)
	if err != nil {
		return false, 0, err
	}

	sub, err := twitchClient.GetSubscriptions(&helix.SubscriptionsParams{
		BroadcasterID: channelId,
		UserID:        []string{userId},
		After:         "",
		First:         0,
	})
	if err != nil {
		return false, 0, err
	}

	if sub.ErrorMessage != "" {
		return false, 0, errors.New(sub.ErrorMessage)
	}

	if len(sub.Data.Subscriptions) == 0 {
		return false, 0, nil
	}

	return true, 0, nil
}

func (c *GiveawaysGrpcImplementation) isFollower(
	ctx context.Context,
	channelId, userId string,
) (bool, *time.Time, error) {
	channel := model.Channels{}
	err := c.db.
		WithContext(ctx).
		Where(`"id" = ?`, channelId).
		First(&channel).
		Error
	if err != nil {
		c.logger.Info("cannot get channel", slog.String("channelId", channelId))
		return false, nil, err
	}

	twitchClient, err := twitch.NewBotClientWithContext(
		ctx,
		channel.BotID,
		c.cfg,
		c.tokensGrpc,
	)
	if err != nil {
		return false, nil, err
	}

	follow, err := twitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: channelId,
			UserID:        userId,
			First:         0,
			After:         "",
		},
	)
	if err != nil {
		return false, nil, err
	}

	if follow.ErrorMessage != "" {
		return false, nil, errors.New(follow.ErrorMessage)
	}

	if len(follow.Data.Channels) == 0 {
		return false, nil, nil
	}

	followedTime := follow.Data.Channels[0].Followed.Time

	return true, &followedTime, nil
}

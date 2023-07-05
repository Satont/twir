package commands_group

import (
	"context"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/commands_group"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommandsGroup struct {
	*impl_deps.Deps
}

func (c *CommandsGroup) convertEntity(g *model.ChannelCommandGroup) *commands_group.Group {
	return &commands_group.Group{
		Id:        &g.ID,
		ChannelId: g.ChannelID,
		Name:      g.Name,
		Color:     g.Color,
	}
}

func (c *CommandsGroup) CommandsGroupsGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*commands_group.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []*model.ChannelCommandGroup
	if err := c.Db.WithContext(ctx).Where(`"channelId" = ?`, dashboardId).Find(&entities).Error; err != nil {
		return nil, err
	}

	return &commands_group.GetAllResponse{
		Groups: lo.Map(entities, func(g *model.ChannelCommandGroup, _ int) *commands_group.Group {
			return c.convertEntity(g)
		}),
	}, nil
}

func (c *CommandsGroup) CommandsGroupsUpdate(
	ctx context.Context,
	req *commands_group.PutRequest,
) (*commands_group.Group, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelCommandGroup{}
	if err := c.Db.WithContext(ctx).Where(`"id" = ? and "channelId" = ?`, req.Id, dashboardId).First(entity).Error; err != nil {
		return nil, err
	}

	entity.Name = req.Group.Name
	entity.Color = req.Group.Color

	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *CommandsGroup) CommandsGroupsCreate(
	ctx context.Context,
	req *commands_group.CreateRequest,
) (*commands_group.Group, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelCommandGroup{
		ID:        uuid.New().String(),
		ChannelID: dashboardId,
		Name:      req.Group.Name,
		Color:     req.Group.Color,
	}
	if err := c.Db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *CommandsGroup) CommandsGroupsDelete(
	ctx context.Context,
	req *commands_group.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	if err := c.Db.
		WithContext(ctx).
		Where(`"id" = ? and "channelId" = ?`, req.Id, dashboardId).
		Delete(&model.ChannelCommandGroup{}).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

package timers

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name            string
	Enabled         bool
	TimeInterval    int
	MessageInterval int
	Responses       []CreateResponse
}

type CreateResponse struct {
	Text       string
	IsAnnounce bool
}

func (c *Service) Create(ctx context.Context, data CreateInput) (entity.Timer, error) {
	createdCount, err := c.timersrepository.CountByChannelID(ctx, data.ChannelID)
	if err != nil {
		return entity.TimerNil, err
	}

	if createdCount >= MaxPerChannel {
		return entity.TimerNil, fmt.Errorf("you can have only %v timers", MaxPerChannel)
	}

	responses := make([]timersrepository.CreateResponse, 0, len(data.Responses))
	for _, response := range data.Responses {
		responses = append(
			responses,
			timersrepository.CreateResponse{
				Text:       response.Text,
				IsAnnounce: response.IsAnnounce,
			},
		)
	}

	timer, err := c.timersrepository.Create(
		ctx,
		timersrepository.CreateInput{
			ChannelID:       data.ChannelID,
			Name:            data.Name,
			Enabled:         data.Enabled,
			TimeInterval:    data.TimeInterval,
			MessageInterval: data.MessageInterval,
			Responses:       responses,
		},
	)

	go func() {
		timersReq := timersbusservice.AddOrRemoveTimerRequest{TimerID: timer.ID.String()}
		if timer.Enabled {
			c.twirbus.Timers.AddTimer.Publish(timersReq)
		} else {
			c.twirbus.Timers.RemoveTimer.Publish(timersReq)
		}
	}()

	c.logger.Audit(
		"Timers create",
		audit.Fields{
			NewValue:      timer,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(timer.ID.String()),
		},
	)

	return c.dbToModel(timer), nil
}
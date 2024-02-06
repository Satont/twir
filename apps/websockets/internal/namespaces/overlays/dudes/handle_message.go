package dudes

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/types"
)

func (c *Dudes) handleMessage(session *melody.Session, msg []byte) {
	userId, ok := session.Get("userId")
	if userId == nil || userId == "" || !ok {
		return
	}

	var overlayId string
	id, ok := session.Get("id")

	if id != nil || ok {
		casted, castOk := id.(string)
		if castOk {
			overlayId = casted
		}
	}

	data := &types.WebSocketMessage{
		CreatedAt: time.Now().UTC().String(),
	}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	if data.EventName == "getSettings" {
		err := c.SendSettings(userId.(string), overlayId)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}
}

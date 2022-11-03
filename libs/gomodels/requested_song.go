package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type RequestedSong struct {
	ID          string    `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	ChannelID   string    `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
	OrderedById string    `gorm:"column:orderedById;type:TEXT;"                   json:"orderedById"`
	VideoID     string    `gorm:"column:videoId;type:varchar;"                    json:"videoId"`
	Title       string    `gorm:"column:title;type:text;"                         json:"title"`
	Duration    int32     `gorm:"column:duration;type:int4"                       json:"duration"`
	CreatedAt   time.Time `gorm:"column:createdAt;type:TIMESTAMP;"                json:"createdAt"`
}

func (c *RequestedSong) TableName() string {
	return "channels_requested_songs"
}

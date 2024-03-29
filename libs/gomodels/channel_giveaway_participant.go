package model

import (
	"github.com/guregu/null"
)

type ChannelGiveawayParticipant struct {
	ID                   string    `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	GiveawayID           string    `gorm:"column:giveaway_id;not null"                                json:"giveaway_id"`
	IsWinner             bool      `gorm:"column:is_winner;not null;default:false"                    json:"is_winner"`
	UserID               string    `gorm:"column:user_id;not null"                                    json:"user_id"`
	DisplayName          string    `gorm:"column:display_name;not null"                               json:"display_name"`
	IsSubscriber         bool      `gorm:"column:is_subscriber;not null;default:false"                json:"is_subscriber"`
	IsFollower           bool      `gorm:"column:is_follower;not null;default:false"                  json:"is_follower"`
	IsModerator          bool      `gorm:"column:is_moderator;not null;default:false"                 json:"is_moderator"`
	IsVip                bool      `gorm:"column:is_vip;not null;default:false"                       json:"is_vip"`
	SubscriberTier       null.Int  `gorm:"column:subscriber_tier"                                     json:"subscriber_tier"`
	UserFollowSince      null.Time `gorm:"column:user_follow_since"                                   json:"user_follow_since"`
	UserStatsWatchedTime int64     `gorm:"column:user_stats_watched_time;not null"                    json:"user_stats_watched_time"`
}

func (ChannelGiveawayParticipant) TableName() string {
	return "channels_giveaways_participants"
}

package model

import (
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
)

type ChannelGiveaway struct {
	ID                        string         `gorm:"primary_key;column:id;type:TEXT;default:uuid_generate_v4()" json:"id"`
	Description               null.String    `gorm:"column:description;type:TEXT;"                              json:"description"`
	ChannelID                 string         `gorm:"column:channel_id;not null"                                 json:"channel_id"`
	CreatedAt                 time.Time      `gorm:"column:created_at;not null;default:now();"                  json:"created_at"`
	FinishedAt                null.Time      `gorm:"column:finished_at;"                                        json:"finished_at"`
	IsRunning                 bool           `gorm:"column:is_running;default:false;not null;"                  json:"is_running"`
	RequiredMinWatchTime      int            `gorm:"column:required_min_watch_time;not null;"                   json:"required_min_watch_time"`
	RequiredMinFollowTime     int            `gorm:"column:required_min_follow_time;not null;"                  json:"required_min_follow_time"`
	RequiredMinMessages       int            `gorm:"column:required_min_messages;not null;"                     json:"required_min_messages"`
	RequiredMinSubscriberTier int            `gorm:"column:required_min_subscriber_tier;not null;"              json:"required_min_subscriber_tier"`
	RequiredMinSubscriberTime int            `gorm:"column:required_min_subscriber_time;not null;"              json:"required_min_subscriber_time"`
	RolesIDS                  pq.StringArray `gorm:"column:roles_ids;type:text[];default:[];"                   json:"roles_ids"`
	Keyword                   string         `gorm:"column:keyword; not null;"                                  json:"keyword"`
	FollowersLuck             int            `gorm:"column:followers_luck;not null;default:1;"                  json:"followers_luck"`
	SubscribersLuck           int            `gorm:"column:subscribers_luck;not null;default:1;"                json:"subscribers_luck"`
	SubscribersTier1Luck      int            `gorm:"column:subscribers_tier1_luck;not null;default:1;"          json:"subscribers_tier1_luck"`
	SubscribersTier2Luck      int            `gorm:"column:subscribers_tier2_luck;not null;default:1;"          json:"subscribers_tier2_luck"`
	SubscribersTier3Luck      int            `gorm:"column:subscribers_tier3_luck;not null;default:1;"          json:"subscribers_tier3_luck"`
	FollowersAgeLuck          bool           `gorm:"column:followers_age_luck;not null;default:false;"          json:"followers_age_luck"`
	WinnersCount              int            `gorm:"column:winners_count;not null;default:1;"                   json:"winners_count"`
	IsFinished                bool           `gorm:"column:is_finished;default:false;not null;"                 json:"is_finished"`
}

func (c *ChannelGiveaway) TableName() string {
	return "channels_giveaways"
}

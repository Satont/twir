-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_giveaways" (
	"id" uuid default uuid_generate_v4() primary key,
	"description" text not null,
	"channel_id" text not null,
	"created_at" timestamp not null default now(),
	"finished_at" timestamp,
	"is_running" boolean not null default 'false',
	"required_min_watch_time" integer not null,
	"required_min_follow_time" integer not null,
	"required_min_messages" integer not null,
	"required_min_subscriber_tier" integer not null,
	"required_min_subscriber_time" integer not null,
	"roles_ids" text[] NOT NULL default '{}',
	"keyword" varchar NOT NULL,
	"followers_luck" integer not null default '0',
	"subscribers_luck" integer not null default '0',
	"subscribers_tier1_luck" integer not null default '0',
	"subscribers_tier2_luck" integer not null default '0',
	"subscribers_tier3_luck" integer not null default '0',
	"followers_age_luck" boolean not null default 'false',
	"winners_count" integer not null default '1',
	"is_finished" boolean not null default 'false'
);

CREATE TABLE "channels_giveaways_participants" (
	"id" uuid default uuid_generate_v4() primary key,
	"giveaway_id" uuid not null,
	"is_winner" boolean not null default false,
	"user_id" text not null,
	"display_name" text not null,
	"is_subscriber" boolean not null default false,
	"is_follower" boolean not null default false,
	"is_moderator" boolean not null default false,
	"is_vip" boolean not null default false,
	"subscriber_tier" integer,
	"user_follow_since" timestamp,
	"user_stats_watched_time" bigint not null
);

ALTER TYPE "channels_roles_permissions_enum" ADD VALUE 'VIEW_GIVEAWAYS';
ALTER TYPE "channels_roles_permissions_enum" ADD VALUE 'MANAGE_GIVEAWAYS';

ALTER TABLE "channels_giveaways" ADD CONSTRAINT "channels_giveaways_channels_channel_fk" FOREIGN KEY ("channel_id") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "channel_giveaways_channel_giveaways_participants_giveaway_fk" FOREIGN KEY ("giveaway_id") REFERENCES "channels_giveaways"("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "channel_giveaways_channel_giveaways_participants_user_fk" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "channel_giveaways_channel_modules_settings_user_fk" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "channels_giveaways_participants_unique" UNIQUE("giveaway_id", "user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
)

func (c *CachedTwitchClient) GetStreamByUserID(
	ctx context.Context,
	userID string,
) (*helix.Stream, error) {
	redisKey := "cache:twir:twitch:streams-by-user-id:" + userID

	if bytes, _ := c.redis.Get(ctx, redisKey).Bytes(); len(bytes) > 0 {
		var helixStream helix.Stream
		if err := json.Unmarshal(bytes, &helixStream); err != nil {
			return nil, err
		}

		return &helixStream, nil
	}

	twitchReq, err := c.client.GetStreams(&helix.StreamsParams{UserIDs: []string{userID}})
	if err != nil {
		return nil, err
	}
	if twitchReq.ErrorMessage != "" {
		return nil, fmt.Errorf("cannot get twitch stream: %s", twitchReq.ErrorMessage)
	}

	if len(twitchReq.Data.Streams) == 0 {
		return nil, fmt.Errorf("stream not found")
	}

	stream := twitchReq.Data.Streams[0]

	streamBytes, err := json.Marshal(stream)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(ctx, redisKey, streamBytes, 50*time.Hour).Err(); err != nil {
		return nil, err
	}

	return &stream, nil
}

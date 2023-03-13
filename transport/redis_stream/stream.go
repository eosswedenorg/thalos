package redis_stream

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Stream struct {
	Key string
	id  string
}

func (s Stream) Read(client *redis.Client, ctx context.Context) error {
	args := &redis.XReadArgs{
		Streams: []string{s.Key, s.id},
	}

	streams, err := client.XRead(ctx, args).Result()
	if err != nil {
		return err
	}

	for _, strm := range streams {
		l := len(strm.Messages)
		if l > 0 {
			id = strm.Messages[l-1].ID

			// Write id to redis
			if err := rs.client.Set(rs.ctx, strm.Stream+":id", rs.id, 0).Err(); err != nil {
				return err
			}
		}
	}
}

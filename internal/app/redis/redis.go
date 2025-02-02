package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"homework/internal/app/pickuppoint"
	"time"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(opt *redis.Options, ttl time.Duration) *Redis {
	return &Redis{
		redis.NewClient(opt),
		ttl,
	}
}

func (r *Redis) GetPointList(ctx context.Context) ([]pickuppoint.PickUpPoint, error) {
	err := r.client.Get(ctx, "points_valid").Err()
	if err != nil {
		return nil, err
	}
	items, err := r.client.LRange(ctx, "points", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	points := make([]pickuppoint.PickUpPoint, len(items))
	for i, item := range items {
		err := json.Unmarshal([]byte(item), &(points[i]))
		if err != nil {
			return nil, err
		}
	}
	return points, nil
}

func (r *Redis) SetPointList(ctx context.Context, points []pickuppoint.PickUpPoint) error {
	err := r.client.Del(ctx, "points").Err()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(points))
	for i, point := range points {
		bytes, err := json.Marshal(point)
		if err != nil {
			return err
		}
		values[i] = string(bytes)
	}
	err = r.client.RPush(ctx, "points", values...).Err()
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "points_valid", true, r.ttl).Err()
}

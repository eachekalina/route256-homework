package cache

import (
	"context"
	"errors"
	"homework/internal/app/pickuppoint"
	"sync"
	"time"
)

type cachedItem[T any] struct {
	value  T
	expire time.Time
}

type Cache struct {
	points            map[uint64]cachedItem[pickuppoint.PickUpPoint]
	pointsMutex       sync.RWMutex
	ttl               time.Duration
	collectorInterval time.Duration
}

func NewCache(ttl time.Duration, collectorInterval time.Duration) *Cache {
	return &Cache{
		points:            make(map[uint64]cachedItem[pickuppoint.PickUpPoint]),
		ttl:               ttl,
		collectorInterval: collectorInterval,
	}
}

func (c *Cache) Run(ctx context.Context) error {
	t := time.NewTicker(c.collectorInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			err := c.invalidateCache(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Cache) invalidateCache(ctx context.Context) error {
	c.pointsMutex.Lock()
	defer c.pointsMutex.Unlock()
	now := time.Now()
	for id, item := range c.points {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if item.expire.Before(now) {
				delete(c.points, id)
			}
		}
	}
	return nil
}

func (c *Cache) PutPoint(point pickuppoint.PickUpPoint) {
	c.pointsMutex.Lock()
	c.points[point.Id] = cachedItem[pickuppoint.PickUpPoint]{
		value:  point,
		expire: time.Now().Add(c.ttl),
	}
	c.pointsMutex.Unlock()
}

func (c *Cache) GetPoint(id uint64) (pickuppoint.PickUpPoint, error) {
	c.pointsMutex.RLock()
	item, ok := c.points[id]
	c.pointsMutex.RUnlock()
	if !ok || item.expire.Before(time.Now()) {
		return pickuppoint.PickUpPoint{}, errors.New("point not found")
	}
	return item.value, nil
}

func (c *Cache) DeletePoint(id uint64) {
	c.pointsMutex.Lock()
	delete(c.points, id)
	c.pointsMutex.Unlock()
}

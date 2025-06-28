package cache

import (
	"order-service/internal/entity"
	"sync"
	"time"

	gcache "github.com/patrickmn/go-cache"
)

type Cache struct {
	mu        sync.Mutex
	cache     *gcache.Cache
	maxSize   int
	orderKeys []string
}

func NewCache(maxSize int, defaultTTL time.Duration) *Cache {
	return &Cache{
		cache:     gcache.New(defaultTTL, 2*defaultTTL),
		maxSize:   maxSize,
		orderKeys: make([]string, 0, maxSize),
	}
}

func (c *Cache) Set(order *entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.orderKeys) >= c.maxSize {
		oldKey := c.orderKeys[0]
		c.cache.Delete(oldKey)
		c.orderKeys = c.orderKeys[1:]
	}

	c.cache.Set(order.OrderUID, order, gcache.DefaultExpiration)
	c.orderKeys = append(c.orderKeys, order.OrderUID)
}

func (c *Cache) Get(orderUID string) (*entity.Order, bool) {
	if val, exists := c.cache.Get(orderUID); exists {
		return val.(*entity.Order), true
	}
	return nil, false
}

func (c *Cache) Load(orders []*entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		if len(c.orderKeys) >= c.maxSize {
			break
		}
		c.cache.Set(order.OrderUID, order, gcache.DefaultExpiration)
		c.orderKeys = append(c.orderKeys, order.OrderUID)
	}
}

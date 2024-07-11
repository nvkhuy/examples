package caching

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	gocache "github.com/patrickmn/go-cache"

	redis "github.com/redis/go-redis/v9"
)

type Client struct {
	cache *cache.ChainCache[string]
}

var client *Client

func New(cnf *config.Configuration) *Client {
	if client != nil {
		return client
	}

	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr: cnf.RedisAddress[0],
	}))

	gocacheStore := gocache_store.NewGoCache(gocache.New(5*time.Minute, 10*time.Minute))

	cacheManager := cache.NewChain[string](
		cache.New[string](redisStore),
		cache.New[string](gocacheStore),
	)

	client = &Client{
		cache: cacheManager,
	}

	return client
}

func GetInstance() *Client {
	if client == nil {
		panic("CacheManager is not initialized")
	}

	return client
}

func (c *Client) Get(key string, data any) error {
	dataStr, err := c.cache.Get(context.Background(), key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(dataStr), data)
}

func (c *Client) Set(key string, data any, options ...store.Option) error {
	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.cache.Set(context.Background(), key, string(byteData), options...)

}

func (c *Client) Delete(key string) error {
	return c.cache.Delete(context.Background(), key)

}

func (c *Client) Clear() error {
	return c.cache.Clear(context.Background())

}

func (c *Client) GetClient() *cache.ChainCache[string] {
	return c.cache
}

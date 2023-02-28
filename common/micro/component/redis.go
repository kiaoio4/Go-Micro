package component

import (
	"context"

	"go-micro/common/micro"
	"go-micro/common/redis"
	goRedis "github.com/go-redis/redis"
)

// TODO RedisComponent should maintain redis client for all dbs, including creating new one if not exist yet

// RedisComponent is Component for redis
type RedisComponent struct {
	micro.EmptyComponent
	client *goRedis.Client
}

// Name of the component
func (c *RedisComponent) Name() string {
	return "Redis"
}

// PreInit called before Init()
func (c *RedisComponent) PreInit(ctx context.Context) error {
	// load config
	redis.SetDefaultRedisConfig()
	return nil
}

// Init the component
func (c *RedisComponent) Init(server *micro.Server) error {
	// init
	//var err error
	redisConf := redis.GetRedisConfig()

	c.client = goRedis.NewClient(&goRedis.Options{
		Addr:     redisConf.Address,
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})
	server.RegisterElement(&micro.RedisElementKey, c.client)
	return nil
}

// PostStop called after Stop()
func (c *RedisComponent) PostStop(ctx context.Context) error {
	// post stop
	return c.client.Close()
}

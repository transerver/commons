package redis

import (
	"crypto/tls"
	"github.com/go-redis/redis"
	"github.com/transerver/commons/configs"
	"github.com/transerver/commons/logger"
	"sync"
)

const Nil = redis.Nil

var client = &redisClient{}

type redisClient struct {
	redis.UniversalClient

	onConnected func(*redis.Conn) error
	tlsConfig   *tls.Config
	config      *configs.RedisConfig
	mutex       sync.Mutex
}

func RegisterOnConnected(fn func(conn *redis.Conn) error) {
	client.onConnected = fn
}

func SetConfig(config *configs.RedisConfig) {
	client.config = config
}

func SetTLSConfig(config *tls.Config) {
	client.tlsConfig = config
}

func Client() *redisClient {
	if client.UniversalClient != nil {
		return client
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if client.UniversalClient != nil {
		return client
	}

	config := client.getConfig()
	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              config.Addrs,
		DB:                 config.DB,
		OnConnect:          client.onConnected,
		Password:           config.Password,
		MaxRetries:         config.MaxRetries,
		MinRetryBackoff:    config.MinRetryBackoff,
		MaxRetryBackoff:    config.MaxRetryBackoff,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
		TLSConfig:          client.tlsConfig,
		MaxRedirects:       config.MaxRedirects,
		ReadOnly:           config.ReadOnly,
		RouteByLatency:     config.RouteByLatency,
		RouteRandomly:      config.RouteRandomly,
		MasterName:         config.MasterName,
	})
	client.UniversalClient = uc
	return client
}

func (c *redisClient) getConfig() *configs.RedisConfig {
	if c.config != nil {
		return c.config
	}

	config, err := configs.FetchRedisConfig()
	if err != nil {
		logger.Panicln(err)
	}
	return config
}

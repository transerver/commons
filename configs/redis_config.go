package configs

import (
	"sync"
	"time"
)

type RedisConfig struct {
	Addrs []string `json:"addrs,omitempty" yaml:"addrs" toml:"addrs"`

	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int `json:"db,omitempty" yaml:"db" toml:"db"`

	// Common options.
	Password           string        `json:"password,omitempty" yaml:"password" toml:"password"`
	MaxRetries         int           `json:"maxRetries,omitempty" yaml:"maxRetries" toml:"maxRetries"`
	MinRetryBackoff    time.Duration `json:"minRetryBackoff,omitempty" yaml:"minRetryBackoff" toml:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `json:"maxRetryBackoff,omitempty" yaml:"maxRetryBackoff" toml:"maxRetryBackoff"`
	DialTimeout        time.Duration `json:"dialTimeout,omitempty" yaml:"dialTimeout" toml:"dialTimeout"`
	ReadTimeout        time.Duration `json:"readTimeout,omitempty" yaml:"readTimeout" toml:"readTimeout"`
	WriteTimeout       time.Duration `json:"writeTimeout,omitempty" yaml:"writeTimeout" toml:"writeTimeout"`
	PoolSize           int           `json:"poolSize,omitempty" yaml:"poolSize" toml:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns,omitempty" yaml:"minIdleConns" toml:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge,omitempty" yaml:"maxConnAge" toml:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout,omitempty" yaml:"poolTimeout" toml:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout,omitempty" yaml:"idleTimeout" toml:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency,omitempty" yaml:"idleCheckFrequency" toml:"idleCheckFrequency"`
	//TLSConfig          *tls.Config   `json:"tlsConfig,omitempty" yaml:"tlsConfig" toml:"tlsConfig"`

	// Only cluster clients.

	MaxRedirects   int  `json:"maxRedirects,omitempty" yaml:"maxRedirects" toml:"maxRedirects"`
	ReadOnly       bool `json:"readOnly,omitempty" yaml:"readOnly" toml:"readOnly"`
	RouteByLatency bool `json:"routeByLatency,omitempty" yaml:"routeByLatency" toml:"routeByLatency"`
	RouteRandomly  bool `json:"routeRandomly,omitempty" yaml:"routeRandomly" toml:"routeRandomly"`

	// The sentinel master name.
	// Only failover clients.
	MasterName string `json:"masterName,omitempty" yaml:"masterName" toml:"masterName"`
}

type redisFetcher struct {
	instance *RedisConfig
	mutex    sync.Mutex
	loaded   bool
}

var rf = &redisFetcher{}

func (f *redisFetcher) Fetch() (*RedisConfig, error) {
	if f.loaded {
		return f.instance, nil
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.loaded {
		return f.instance, nil
	}

	err := Fetch(redisKey, &f.instance)
	if err != nil {
		return nil, err
	}
	f.loaded = true
	return f.instance, nil
}

func FetchRedisConfig() (*RedisConfig, error) {
	return rf.Fetch()
}

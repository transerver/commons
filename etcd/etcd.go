package etcd

import (
	"github.com/gookit/color"
	"github.com/transerver/commons/logger"
	v3 "go.etcd.io/etcd/client/v3"
	"os"
	"strings"
	"sync"
	"time"
)

type etcdClient struct {
	*v3.Client
	config      *v3.Config
	onConnected func(*v3.Client)
	endpoints   []string
	mutex       sync.Mutex
}

var ec = new(etcdClient)

type Option func(*etcdClient)

func WithConfig(config v3.Config) Option {
	return func(c *etcdClient) {
		c.config = &config
	}
}

// WithOnConnected executed on etcd connect successfully
func WithOnConnected(onConnected func(client *v3.Client)) Option {
	return func(c *etcdClient) {
		c.onConnected = onConnected
	}
}

// WithEndpoints when v3.Client not have endpoints
// then override the default endpoints
// otherwise can be used default endpoints
func WithEndpoints(endpoints []string) Option {
	return func(c *etcdClient) {
		c.endpoints = endpoints
	}
}

func NewClient(ops ...Option) *v3.Client {
	c := new(etcdClient)
	for _, op := range ops {
		op(c)
	}
	c.connect()
	return c.Client
}

func OnConnected(fn func(*v3.Client)) {
	WithOnConnected(fn)(ec)
}

func RegisterConfig(config v3.Config) {
	WithConfig(config)(ec)
}

// SetEndpoints can be overridden RegisterConfig's config endpoints
func SetEndpoints(endpoints []string) {
	WithEndpoints(endpoints)(ec)
}

// Client returns the default etcd client
// there's no config and use default endpoint
func Client() *v3.Client {
	if ec.Client == nil {
		ec.connect()
	}
	return ec.Client
}

func (c *etcdClient) connect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Client != nil {
		return
	}

	var err error
	if c.config == nil {
		c.config = &v3.Config{DialTimeout: 30 * time.Second}
	}

	c.config.Endpoints = c.getEndpoints()
	c.Client, err = v3.New(*c.config)
	if err != nil {
		logger.Panicln("initialize etcd error:", err)
	}

	if c.onConnected != nil {
		c.onConnected(c.Client)
	}
	return
}

func (c *etcdClient) getEndpoints() []string {
	if len(c.endpoints) > 0 {
		return c.endpoints
	}

	if len(c.config.Endpoints) > 0 {
		return c.config.Endpoints
	}

	endpoints := os.Getenv("ETCD_ENDPOINTS")
	if len(endpoints) == 0 {
		logger.Warnf(color.HiYellow.Sprintf("can't find etcd endpoints in environment, use default address[%s]", endpoints))
		return []string{"127.0.0.1:2379"}
	}
	return strings.Split(endpoints, ",")
}

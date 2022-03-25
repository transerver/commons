package configs

import (
	"bytes"
	"github.com/BurntSushi/toml"
	json "github.com/json-iterator/go"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/stretchr/testify/require"
	"github.com/transerver/commons/logger"
	"github.com/xo/dburl"
	"testing"
	"time"
)

func TestRedisConfig(t *testing.T) {
	config := RedisConfig{
		Addrs:              []string{":6379"},
		DB:                 10,
		Password:           "Password",
		MaxRetries:         10,
		MinRetryBackoff:    time.Second,
		MaxRetryBackoff:    time.Second,
		DialTimeout:        time.Second,
		ReadTimeout:        time.Second,
		WriteTimeout:       time.Second,
		PoolSize:           10,
		MinIdleConns:       10,
		MaxConnAge:         time.Second,
		PoolTimeout:        time.Second,
		IdleTimeout:        time.Second,
		IdleCheckFrequency: time.Second,
		MaxRedirects:       10,
		ReadOnly:           true,
		RouteByLatency:     true,
		RouteRandomly:      true,
		MasterName:         "MasterName",
	}
	data, err := json.Marshal(config)
	require.NoError(t, err)
	t.Logf("JSON: %s", data)

	var w bytes.Buffer
	encoder := toml.NewEncoder(&w)
	err = encoder.Encode(config)
	require.NoError(t, err)
	t.Logf("TOML: %s", w.String())
}

func TestLoadRedisConfigFromRemote(t *testing.T) {
	err := AddEtcdProvider("/configuration/configs/config.json")
	require.NoError(t, err)
	err = ReadConfig()
	require.NoError(t, err)

	var config RedisConfig
	err = Fetch("redis", &config)
	require.NoError(t, err)
	require.NotEmpty(t, config)
}

func TestLoadRedisConfigFromFile(t *testing.T) {
	SetConfigFile("./config.toml")
	err := ReadConfig()
	require.NoError(t, err)

	var config RedisConfig
	err = Fetch("redis", &config)
	require.NoError(t, err)
	require.NotEmpty(t, config)
}

func TestViperRemote(t *testing.T) {
	var err error
	viper.RemoteConfig = &etcdConfigFactory{}
	err = viper.AddRemoteProvider("etcd", "127.0.0.1:2379", "config.json")
	require.NoError(t, err)
	viper.SetConfigType("json")
	err = viper.ReadRemoteConfig()
	require.NoError(t, err)

	var databases []*DBConfig
	err = viper.UnmarshalKey("databases", &databases)
	require.NoError(t, err)
}

func TestFetch(t *testing.T) {
	err := AddEtcdProvider("/configuration/configs/config.json")
	require.NoError(t, err)
	err = ReadConfig()
	require.NoError(t, err)
	var dbs []*DBConfig
	dbs, err = FetchDBConfigs()
	require.NoError(t, err)
	require.NotEmpty(t, dbs)

	var address string
	err = Fetch("port", &address)
	require.NoError(t, err)
	require.NotEmpty(t, address)

	config, err := FetchRedisConfig()
	require.NoError(t, err)
	require.NotEmpty(t, config)
}

func TestAsyncFetch(t *testing.T) {
	err := AddEtcdProvider("/configuration/configs/config.json")
	require.NoError(t, err)
	err = ReadConfig()
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		go func(i int) {
			logger.Infof("第%d次拉取", i)
			dbs, err := FetchDBConfigs()
			require.NoError(t, err)
			require.NotEmpty(t, dbs)
		}(i)
	}

	<-make(chan struct{})
}

func TestParseURL(t *testing.T) {
	urlstr := "postgres://charlie:root@127.0.0.1:54321/configuration?sslmode=disable"
	url, err := dburl.Parse(urlstr)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	t.Logf("URL: %+v, DSN: %s", url, url.DSN)

	postgres, err := dburl.GenPostgres(url)
	require.NoError(t, err)
	require.NotEmpty(t, postgres)

	url, dbName, baseDSN, err := ParseURL(urlstr)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.NotEmpty(t, dbName)
	require.NotEmpty(t, baseDSN)
	t.Logf("URL: %s, dbName: %s, baseDSN: %s", url, dbName, baseDSN)
}

func TestViperWithFile(t *testing.T) {
	viper.SetConfigFile("./config.toml")
	err := viper.ReadInConfig()
	require.NoError(t, err)

	var dbs []DBConfig
	err = viper.UnmarshalKey("databases", &dbs)
	require.NoError(t, err)

	var redis []RedisConfig
	err = viper.UnmarshalKey("redis", &redis)
	require.NoError(t, err)
}

func TestViperFile(t *testing.T) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	require.NoError(t, err)

	settings := viper.AllSettings()
	var dbs []DBConfig
	err = viper.UnmarshalKey("databases", &dbs)
	require.NoError(t, err)
	t.Logf("%+v", settings)

	var redis []RedisConfig
	err = viper.UnmarshalKey("redis", &redis)
	require.NoError(t, err)
}

func TestGetDatabaseConfig(t *testing.T) {
	config, err := FetchDBConfigs()
	require.NoError(t, err)
	require.NotEmpty(t, config)
	config, err = FetchDBConfigs()
	require.NoError(t, err)
	require.NotEmpty(t, config)
	config, err = FetchDBConfigs()
	require.NoError(t, err)
	require.NotEmpty(t, config)

	dbConfig := FetchDBConfigWithName("configuration")
	require.NotEmpty(t, dbConfig)
}

func TestClearCode(t *testing.T) {
	dbConfig := FetchDBConfigWithName("configuration")
	require.NotEmpty(t, dbConfig)
}

package configs

import (
	"fmt"
	e2 "github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/transerver/commons/etcd"
	"github.com/xo/dburl"
	"path/filepath"
	"strings"
)

const (
	databaseKey = "databases"
	redisKey    = "redis"
)

type ConfigNotFoundErr struct {
	err error
}

func (c ConfigNotFoundErr) Error() string {
	return fmt.Sprintf("No configuration found -> %s", c.err.Error())
}

func AddEtcdProvider(paths ...string) (err error) {
	viper.RemoteConfig = &etcdConfigFactory{}
	viper.SetConfigType("json")
	path := filepath.Join(paths...)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	err = viper.AddRemoteProvider("etcd", etcd.Client().Endpoints()[0], path)
	return
}

func SetConfigFile(path string) {
	viper.SetConfigFile(path)
}

func ReadConfig() error {
	var bothNotFound bool
	rerr := viper.ReadRemoteConfig()
	if rerr != nil {
		_, ok := rerr.(viper.RemoteConfigError)
		if !ok {
			return rerr
		}
		bothNotFound = ok
	}

	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return err
		}
		bothNotFound = bothNotFound && ok
		if rerr != nil {
			err = e2.Wrap(err, rerr.Error())
		}
	} else {
		bothNotFound = false
	}
	if bothNotFound {
		return ConfigNotFoundErr{err}
	}
	return nil
}

func Fetch(key string, v interface{}, opts ...viper.DecoderConfigOption) (err error) {
	return viper.UnmarshalKey(key, v, opts...)
}

// ParseURL parse the url to dburl.URL, dbName baseDSN(without password)
// error on parse fail
func ParseURL(configURL string) (url *dburl.URL, dbName, baseDSN string, err error) {
	url, err = dburl.Parse(configURL)
	if err != nil {
		return
	}

	dbName = url.Path[1:]
	if pwd, ok := url.User.Password(); ok {
		baseDSN = strings.Replace(url.DSN, "password="+pwd, "password=[***]", 1)
	} else {
		baseDSN = url.DSN
	}
	return
}

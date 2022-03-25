package configs

import (
	"github.com/transerver/commons/colors"
	"github.com/transerver/commons/logger"
	"sync"
	"time"
)

type DBConfig struct {
	Driver string `json:"driver,omitempty" toml:"driver" yaml:"driver"`
	DBName string `json:"dbName,omitempty" toml:"dbName" yaml:"dbName"`
	DSN    string `json:"dsn,omitempty" toml:"dsn" yaml:"dsn"`

	DesensitiseDSN string `json:"desensitiseDsn,omitempty" toml:"desensitiseDsn" yaml:"desensitiseDsn"`

	URL     string `json:"url,omitempty" toml:"url" yaml:"url"`
	Options struct {
		MaxOpenConns    int           `json:"maxOpenConns,omitempty" toml:"maxOpenConns" yaml:"maxOpenConns"`
		MaxIdleConns    int           `json:"maxIdleConns,omitempty" toml:"maxIdleConns" yaml:"maxIdleConns"`
		ConnMaxIdleTime time.Duration `json:"connMaxIdleTime,omitempty" toml:"connMaxIdleTime" yaml:"connMaxIdleTime"`
		ConnMaxLifetime time.Duration `json:"connMaxLifetime,omitempty" toml:"connMaxLifeTime" yaml:"connMaxLifeTime"`
	} `json:"options,omitempty" toml:"options" yaml:"options"`
}

type dbFetcher struct {
	instance []*DBConfig
	mutex    sync.Mutex
	loaded   bool
}

var dbf = &dbFetcher{}

// dbHandler same db name be override, keep the last one
func dbHandler(value interface{}) {
	configs, ok := value.(*[]*DBConfig)
	if !ok || configs == nil {
		return
	}

	names := make(map[string]struct{})
	for i, config := range *configs {
		url, dbName, baseDSN, err := ParseURL(config.URL)
		if err != nil {
			logger.Errorf("parse database url fail: %+v", err)
			continue
		}
		config.DesensitiseDSN = baseDSN
		config.DBName = dbName
		config.DSN = url.DSN
		config.Driver = url.Driver

		l := len(names)
		names[dbName] = struct{}{}
		if len(names) == l {
			pre := (*configs)[i-1]
			logger.Warn(colors.HiYellowUnderline.Sprintf("database name with (%s) is duplicate, will be ignore.", pre.DesensitiseDSN))
			*configs = append((*configs)[:i-1], (*configs)[i:]...)
		}
	}
	return
}

func (f *dbFetcher) FetchConfigs() ([]*DBConfig, error) {
	if f.loaded {
		return f.instance, nil
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.loaded {
		return f.instance, nil
	}

	err := Fetch(databaseKey, &f.instance)
	if err != nil {
		logger.Warn(colors.HiYellowUnderline.Sprint(err.Error()))
		return nil, err
	}
	dbHandler(&f.instance)
	f.loaded = true
	return f.instance, nil
}

func FetchDBConfigs() ([]*DBConfig, error) {
	return dbf.FetchConfigs()
}

func FetchDBConfigWithName(dbName string) *DBConfig {
	configs, err := FetchDBConfigs()
	if err != nil {
		return nil
	}
	var rf *DBConfig
	for _, config := range configs {
		if config.DBName == dbName {
			rf = config
			break
		}
	}
	return rf
}

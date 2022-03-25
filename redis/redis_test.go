package redis

import (
	"github.com/stretchr/testify/require"
	"github.com/transerver/commons/configs"
	"github.com/transerver/commons/logger"
	"testing"
	"time"
)

func init() {
	configs.SetConfigFile("../configs/config.toml")
	_ = configs.ReadConfig()
}

func TestLock(t *testing.T) {
	lock, err := Obtain(time.Second*5, "test %s", "obtain")
	if err != nil || !lock.Locked {
		return
	}
	defer lock.LoggedRelease()

	logger.Info("start obtain")
	<-time.After(time.Second * 2)
	ttl, err := lock.TTL()
	require.NoError(t, err)
	logger.Info("before refresh TTL:", ttl.String())
	err = lock.Refresh(time.Second * 10)
	require.NoError(t, err)
	ttl, err = lock.TTL()
	require.NoError(t, err)
	logger.Info("after refresh TTL:", ttl.String())

	logger.Infof("end obtain, Key: %s, Token: %s", lock.Key, lock.value)
}

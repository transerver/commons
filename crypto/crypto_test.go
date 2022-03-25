package crypto

import (
	"github.com/stretchr/testify/require"
	"github.com/transerver/commons/configs"
	"github.com/transerver/commons/utils"
	"testing"
)

func TestFetchRSAKey(t *testing.T) {
	configs.SetConfigFile("../configs/config.toml")
	err := configs.ReadConfig()
	require.NoError(t, err)

	requestId := utils.RandomString(5)

	key, err := FetchRsaKey(requestId, WithBits(20))
	require.NoError(t, err)
	require.NotEmpty(t, key)
	key, err = FetchRsaKey(requestId)
	require.NoError(t, err)
	require.NotEmpty(t, key)
	key, err = FetchRsaKey("ptgwh")
	require.NoError(t, err)
	require.NotEmpty(t, key)
}

package configs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/transerver/commons/etcd"
	"io"
)

type etcdConfigFactory struct{}

func (f *etcdConfigFactory) Get(rp viper.RemoteProvider) (io.Reader, error) {
	ctx := context.TODO()
	response, err := etcd.Client().Get(ctx, rp.Path())
	if err != nil {
		return nil, err
	}

	if response.Count == 0 {
		return nil, fmt.Errorf("configuration item not found in etcd")
	}

	data := response.Kvs[response.Count-1].Value
	return bytes.NewReader(data), err
}

func (f *etcdConfigFactory) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return f.Get(rp)
}

func (f *etcdConfigFactory) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	return nil, nil
}

package etcd

import (
	"context"
	"github.com/ohwin/micro/core/config"
	"github.com/ohwin/micro/core/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var Client *clientv3.Client

func Init() {
	c := config.App.Etcd
	if c.IsNil() {
		return
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.App.Etcd.Endpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Error("init etcd config error, ", err)
		panic(err)
	}

	for _, endpoint := range cli.Endpoints() {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		_, err = cli.Status(ctx, endpoint)
		if err != nil {
			cancel()
			log.Errorf("connection etcd(%s) failure: %s", endpoint, err)
			panic("connection etcd failure")
		}
		cancel()
	}

	Client = cli
}

package etcd

import (
	"context"
	"github.com/ohwin/micro/core/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Cluster struct {
	Ctx   context.Context
	Lease *clientv3.LeaseGrantResponse `json:"lease"`
	Key   string                       `json:"key"`
	Value string                       `json:"value"`
}

const TimeToLive = 10

type Option func(cluster *Cluster)

var Nodes map[string]*Cluster = make(map[string]*Cluster, 0)

func NewCluster(key string, value string, opt ...Option) *Cluster {
	cluster := &Cluster{
		Ctx:   context.TODO(),
		Key:   key,
		Value: value,
	}
	for _, f := range opt {
		f(cluster)
	}
	Nodes[key] = cluster
	return cluster
}

// Register 注册服务到etcd
func (c *Cluster) Register() {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	_, err := Client.Put(ctx, c.Key, c.Value, clientv3.WithLease(c.Lease.ID))
	if err != nil {
		log.Errorf("register server to etcd timeout")
		panic(err)
	}
}

func (c *Cluster) KeepLive() {
	ch, err := Client.KeepAlive(c.Ctx, c.Lease.ID)
	if err != nil {
		panic(err)
		return
	}

	go func() {
		for {
			select {
			case x := <-ch:
				if x == nil {
					log.Error("etcd connect failure, try reconnecting...")
					time.Sleep(2 * time.Second)
					c.Reload()
					return
				} else {
					log.Info(x)
				}
			}
		}
	}()
}

// Reload 服务重连
func (c *Cluster) Reload() {
	Lease(c)
	c.Register()
	c.KeepLive()
	c.Watch()
}

func (c *Cluster) Watch() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		watchChan := Client.Watch(ctx, c.Key, clientv3.WithPrefix())
		for watch := range watchChan {
			for _, event := range watch.Events {
				switch event.Type {
				case mvccpb.PUT:
					Nodes[string(event.Kv.Key)] = NewCluster(string(event.Kv.Key), string(event.Kv.Value))
					log.Errorf("etcd(%s) online", string(event.Kv.Key))
				case mvccpb.DELETE:
					delete(Nodes, string(event.Kv.Key))
					log.Errorf("etcd(%s) offline", string(event.Kv.Key))
				}
			}
		}

	}()
}

// WithLease 设置租约
func WithLease(ttl int64) Option {

	return func(cluster *Cluster) {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		lease, err := Client.Grant(ctx, ttl)
		if err != nil {
			panic(err)
		}

		cluster.Lease = lease
	}
}

func Lease(cluster *Cluster) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	lease, err := Client.Grant(ctx, cluster.Lease.TTL)
	if err != nil {
		panic(err)
	}

	cluster.Lease = lease

}

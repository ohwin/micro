package resolver

import (
	"context"
	"github.com/ohwin/micro/core/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"strings"
)

type EtcdResolver struct {
	target    resolver.Target
	cc        resolver.ClientConn
	addresses []resolver.Address
}

func (r *EtcdResolver) ResolveNow(options resolver.ResolveNowOptions) {

	err := r.cc.UpdateState(resolver.State{
		Addresses: r.addresses,
	})
	if err != nil {
		panic(err)
		return
	}
}

func (r *EtcdResolver) Close() {

}

type EtcdResolverBuilder struct {
}

func (b *EtcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	key := strings.Trim(target.URL.Path, "/")
	resp, err := etcd.Client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	addresses := make([]resolver.Address, 0)
	for _, kv := range resp.Kvs {
		addresses = append(addresses, resolver.Address{Addr: string(kv.Value)})
	}

	r := &EtcdResolver{
		target:    target,
		cc:        cc,
		addresses: addresses,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (b *EtcdResolverBuilder) Scheme() string {
	return "etcd"
}

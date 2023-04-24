package rpc

import (
	"fmt"
	"github.com/ohwin/micro/core/config"
	"github.com/ohwin/micro/core/etcd"
	"github.com/ohwin/micro/core/initialize"
	"github.com/ohwin/micro/core/rpc/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

func init() {
	resolver.Register(&resolver.EtcdResolverBuilder{})
}

type Server struct {
	Listen net.Listener
	Server *grpc.Server
}

func NewServer() *Server {
	initialize.Init()
	listen, err := net.Listen("tcp", config.Addr())
	if err != nil {
		panic(err)
		return nil
	}
	return &Server{
		Listen: listen,
		Server: grpc.NewServer(),
	}
}

// Register 注册服务
func (s *Server) Register(desc *grpc.ServiceDesc, srv interface{}) *Server {

	s.Server.RegisterService(desc, srv)

	key := fmt.Sprintf("%s/%s", desc.ServiceName, config.Addr())
	cluster := etcd.NewCluster(key, config.Addr(), etcd.WithLease(etcd.TimeToLive))
	cluster.Register()
	cluster.KeepLive()
	cluster.Watch()
	return s
}

func (s *Server) Serve() {
	err := s.Server.Serve(s.Listen)
	if err != nil {
		panic(err)
	}
}

// Client 获取服务连接
func Client(desc *grpc.ServiceDesc, f func(conn *grpc.ClientConn) interface{}) interface{} {
	target := fmt.Sprintf("etcd://%s/%s", config.App.Etcd.Endpoints[0], desc.ServiceName)
	conn, err := grpc.Dial(target,
		NewCredentials(),
		RoundRobin())
	if err != nil {
		panic(err)
	}
	return f(conn)
}

// RoundRobin 轮询
func RoundRobin() grpc.DialOption {
	return grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)
}

func NewCredentials() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

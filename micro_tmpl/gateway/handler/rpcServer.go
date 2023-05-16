package handler

import (
	"config/configHandle"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"micro_tmpl/gateway/proto"
	"sync"
)

var rpcServerOnce sync.Once
var rpcInstance *RpcServer

type RpcServer struct {
}

func GetRpcServer() *RpcServer {
	rpcServerOnce.Do(func() {
		rpcInstance = &RpcServer{}
	})
	return rpcInstance
}
func (s *RpcServer) getPlayerConn() (proto.PlayerClient, *grpc.ClientConn) {
	generalConfig := configHandle.GetConfigManInstance().GetGeneralConfig()
	etcdUrl := generalConfig.EtcdConfig.EtcdUrl
	serverName := generalConfig.EtcdConfig.EtcdServerName
	cli, err := clientv3.NewFromURL(etcdUrl)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	grpcEtcdUrl := fmt.Sprintf("etcd:///%s", serverName)
	conn, err := grpc.Dial(grpcEtcdUrl, grpc.WithResolvers(etcdResolver), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := proto.NewPlayerClient(conn)
	return c, conn
}

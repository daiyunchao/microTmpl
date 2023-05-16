package handler

import (
	"common/lib"
	"config/configHandle"
	"context"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"log"
	"micro_tmpl/player/playerModel"
	"micro_tmpl/player/proto"
	"micro_tmpl/player/service"
	"net"
)

type ProtoServer struct {
	proto.UnimplementedPlayerServer
	Mongo *lib.MongoStore
}

func (s *ProtoServer) GetInfoById(ctx context.Context, req *proto.ReqGetInfoById) (*proto.ResGetInfo, error) {
	service := service.PlayerService{}
	service.SetMongo(s.Mongo)
	return nil, nil
}
func (s *ProtoServer) GetInfoByName(ctx context.Context, req *proto.ReqGetInfoByName) (*proto.ResGetInfo, error) {
	service := service.PlayerService{}
	service.SetMongo(s.Mongo)
	playerInfo, err := service.GetPlayerInfoByName(ctx, req.Name)
	if err != nil {
		return &proto.ResGetInfo{}, err
	}
	if playerInfo == nil {
		return &proto.ResGetInfo{}, nil
	}
	res := &proto.ResGetInfo{
		Id:       playerInfo.Id,
		Name:     playerInfo.Name,
		Password: playerInfo.Password,
	}
	return res, nil
}

func (s *ProtoServer) Register(ctx context.Context, req *proto.ReqRegister) (*proto.ResGetInfo, error) {
	service := service.PlayerService{}
	service.SetMongo(s.Mongo)
	playerInfo := &playerModel.Player{
		Name:     req.Name,
		Password: req.Password,
	}
	err := service.InsertPlayer(ctx, playerInfo)
	if err != nil {
		return nil, err
	}
	resPlayer, err := service.GetPlayerInfoByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if playerInfo == nil {
		return &proto.ResGetInfo{}, nil
	}
	res := &proto.ResGetInfo{
		Id:       resPlayer.Id,
		Name:     resPlayer.Name,
		Password: resPlayer.Password,
	}
	return res, nil
}
func (s *ProtoServer) RegisterRpcServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	rpcServer := grpc.NewServer()
	proto.RegisterPlayerServer(rpcServer, &ProtoServer{
		Mongo: s.Mongo,
	})
	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *ProtoServer) RegisterEtcdServer(address string) error {
	serverConfig := configHandle.GetConfigManInstance().GetGeneralConfig()
	etcdUrl := serverConfig.EtcdConfig.EtcdUrl
	serverName := serverConfig.EtcdConfig.EtcdServerName
	var ttl int64 = 100
	etcdClient, err := clientV3.NewFromURL(etcdUrl)
	if err != nil {
		return err
	}
	em, err := endpoints.NewManager(etcdClient, serverName)
	if err != nil {
		return err
	}
	lease, _ := etcdClient.Grant(context.TODO(), ttl)
	err = em.AddEndpoint(context.TODO(), fmt.Sprintf("%s/%s", serverName, address), endpoints.Endpoint{Addr: address}, clientV3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	etcdClient.KeepAlive(context.TODO(), lease.ID)
	return err
}

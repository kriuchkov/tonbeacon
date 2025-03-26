package grpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

type TonBeacon struct {
	pb.UnimplementedTonBeaconServer
	accountSvc ports.AccountServicePort
	server     *grpc.Server
}

func NewTonBeacon(account ports.AccountServicePort) *TonBeacon {
	return &TonBeacon{accountSvc: account, server: grpc.NewServer()}
}

func (s *TonBeacon) Run(lis net.Listener) error {
	pb.RegisterTonBeaconServer(s.server, s)
	reflection.Register(s.server)
	return s.server.Serve(lis)
}

func (s *TonBeacon) Stop() {
	s.server.Stop()
}

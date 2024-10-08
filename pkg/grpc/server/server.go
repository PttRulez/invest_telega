package server

import (
	"context"

	tgGrpc "github.com/pttrulez/invest_telega/pkg/grpc"
)

func (s *GRPCTelegaServer) SendMsg(ctx context.Context, msgInfo *tgGrpc.MessageInfo) (*tgGrpc.None, error) {
	return &tgGrpc.None{}, s.svc.SendMsg(msgInfo)
}

type GRPCTelegaServer struct {
	tgGrpc.UnimplementedTelegaServer
	svc Telega
}

type Telega interface {
	SendMsg(msgInfo *tgGrpc.MessageInfo) error
}

func NewGRPCTelegaServer(svc Telega) *GRPCTelegaServer {
	return &GRPCTelegaServer{
		svc: svc,
	}
}

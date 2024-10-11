package grpctransport

import (
	"context"

	"github.com/pttrulez/invest_telega/pkg/protogen"
)

func (s *GRPCTelegaServer) SendMsg(ctx context.Context, msgInfo *protogen.MessageInfo) (
	*protogen.None, error) {
	return &protogen.None{}, s.svc.SendMsg(msgInfo)
}

type GRPCTelegaServer struct {
	protogen.UnimplementedTelegaServer
	svc Telega
}

type Telega interface {
	SendMsg(msgInfo *protogen.MessageInfo) error
}

func NewGRPCTelegaServer(svc Telega) *GRPCTelegaServer {
	return &GRPCTelegaServer{
		svc: svc,
	}
}

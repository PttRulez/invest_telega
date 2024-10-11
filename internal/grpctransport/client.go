package grpctransport

import (
	"context"

	"github.com/pttrulez/investor-go-next/go-api/pkg/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCInvestorClient struct {
	protogen.InvestorClient
	grpcClient protogen.InvestorClient
}

func NewInvestorGRPCClient(endpoint string) (*GRPCInvestorClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	grpcClient := protogen.NewInvestorClient(conn)

	return &GRPCInvestorClient{grpcClient: grpcClient}, nil
}

func (s *GRPCInvestorClient) GetPortfolioList(ctx context.Context,
	chatId string) ([]*protogen.Portfolio, error) {
	req := &protogen.PortfolioListRequest{
		ChatId: chatId,
	}
	res, err := s.grpcClient.GetPortfolioList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.GetPortfolios(), nil
}

func (s *GRPCInvestorClient) GetPortfolioSummaryMessage(ctx context.Context, portfolioID int,
	chatId string) (string, error) {
	req := &protogen.PortfolioRequest{
		Id:     int64(portfolioID),
		ChatId: chatId,
	}
	res, err := s.grpcClient.GetPortfolioSummaryMessage(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetText(), nil
}

package grpctransport

import (
	"context"

	"github.com/pttrulez/invest_telega/pkg/protogen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TelegramClient struct {
	protogen.TelegaClient
	client protogen.TelegaClient
}

func NewTelegramClient(endpoint string) (*TelegramClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := protogen.NewTelegaClient(conn)

	return &TelegramClient{client: c}, nil
}

func (c *TelegramClient) SendMsg(ctx context.Context, msgInfo *protogen.MessageInfo) error {
	_, err := c.client.SendMsg(ctx, msgInfo)
	return err
}

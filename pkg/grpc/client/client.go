package client

import (
	"context"
	tgGrpc "github.com/pttrulez/invest_telega/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TelegramClient struct {
	tgGrpc.TelegaClient
	client tgGrpc.TelegaClient
}

func NewTelegramClient(endpoint string) (*TelegramClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := tgGrpc.NewTelegaClient(conn)

	return &TelegramClient{client: c}, nil
}

func (c *TelegramClient) SendMsg(ctx context.Context, msgInfo *tgGrpc.MessageInfo) error {
	_, err := c.client.SendMsg(ctx, msgInfo)
	return err
}

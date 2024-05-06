package sender

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/sender"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(ctx context.Context, addr string, log *slog.Logger) (*Client, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sender service: %w", err)
	}

	return &Client{
		api: pb.NewSenderServiceClient(cc),
	}, nil
}

type Client struct {
	api pb.SenderServiceClient
	log *slog.Logger
}

func (c *Client) Send(ctx context.Context, msg string) error {
	_, err := c.api.Send(ctx, &pb.SendRequest{Message: msg})
	if err != nil {
		return fmt.Errorf("failed to create stream to the sender: %w", err)
	}
	return nil
}

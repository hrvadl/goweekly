package translator

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/translator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(ctx context.Context, addr string, log *slog.Logger) (*Client, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to translator service: %w", err)
	}

	return &Client{
		api: pb.NewTranslateServiceClient(cc),
	}, nil
}

type Client struct {
	api pb.TranslateServiceClient
	log *slog.Logger
}

func (c *Client) Translate(ctx context.Context, msg string) (string, error) {
	res, err := c.api.Translate(ctx, &pb.TranslateRequest{Message: msg})
	if err != nil {
		return "", fmt.Errorf("failed to translate the message: %w", err)
	}

	return res.Message, nil
}

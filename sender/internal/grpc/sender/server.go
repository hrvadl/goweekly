package sender

import (
	"errors"
	"fmt"
	"io"

	"github.com/hrvadl/goweeky/protos/gen/go/v1/sender"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hrvadl/goweeky/sender/internal/platform"
)

func Register(gRPC *grpc.Server, s Sender) {
	sender.RegisterSenderServiceServer(gRPC, &server{
		sender: s,
	})
}

type Sender interface {
	Send(ctx context.Context, msg platform.Message) error
}

type server struct {
	sender.UnimplementedSenderServiceServer
	sender Sender
}

func (srv *server) Send(s sender.SenderService_SendServer) error {
	ctx := s.Context()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("time out reached")
		default:
			msg, err := s.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}

			err = srv.sender.Send(ctx, platform.Message{
				Message:   msg.Message,
				ChatID:    msg.ChatId,
				ParseMode: msg.ParseMode,
			})
			if err != nil {
				return err
			}
		}
	}
}

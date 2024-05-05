package sender

import (
	"errors"
	"fmt"
	"io"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/sender"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hrvadl/goweekly/sender/internal/platform/sender"
)

func Register(gRPC *grpc.Server, s Sender) {
	pb.RegisterSenderServiceServer(gRPC, &server{
		sender: s,
	})
}

type Sender interface {
	Send(ctx context.Context, msg sender.Message) error
}

type server struct {
	pb.UnimplementedSenderServiceServer
	sender Sender
}

func (srv *server) Send(s pb.SenderService_SendServer) error {
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

			err = srv.sender.Send(ctx, sender.Message{
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

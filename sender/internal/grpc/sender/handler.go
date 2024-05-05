package sender

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/sender"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hrvadl/goweekly/sender/internal/platform/sender"
)

func Register(gRPC *grpc.Server, s Sender, l *slog.Logger) {
	pb.RegisterSenderServiceServer(gRPC, &server{
		sender: s,
		log:    l,
	})
}

type Sender interface {
	Send(ctx context.Context, msg sender.Message) error
}

type server struct {
	pb.UnimplementedSenderServiceServer
	sender Sender
	log    *slog.Logger
}

func (srv *server) Send(s pb.SenderService_SendServer) error {
	srv.log.Info("Got a send streaming request")
	doneCh := make(chan struct{})
	errCh := make(chan error)
	ctx := s.Context()

	go srv.handleStream(s, doneCh, errCh)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("time out reached")
		case err := <-errCh:
			return fmt.Errorf("failed to send msg: %w", err)
		case <-doneCh:
			return nil
		}
	}
}

func (srv *server) handleStream(
	s pb.SenderService_SendServer,
	doneCh chan<- struct{},
	errCh chan<- error,
) {
	breakCh := make(chan struct{})
	ctx := s.Context()

loop:
	for {
		msg, err := s.Recv()
		if errors.Is(err, io.EOF) {
			doneCh <- struct{}{}
			return
		}

		select {
		case <-breakCh:
			break loop
		default:
			go func() {
				err = srv.sender.Send(ctx, sender.Message{
					Message: msg.Message,
				})
				if err != nil {
					breakCh <- struct{}{}
					errCh <- err
				}
			}()
		}
	}
}

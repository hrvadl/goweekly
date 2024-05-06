package sender

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/sender"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hrvadl/goweekly/sender/internal/platform/sender"
)

func Register(gRPC *grpc.Server, s Sender, l *slog.Logger) {
	pb.RegisterSenderServiceServer(gRPC, &server{
		sender: s,
		log:    l,
		donech: make(chan struct{}),
		errch:  make(chan error),
		msgch:  make(chan *pb.SendRequest),
	})
}

type Sender interface {
	Send(ctx context.Context, msg sender.Message) error
}

type server struct {
	pb.UnimplementedSenderServiceServer
	sender Sender
	log    *slog.Logger
	donech chan struct{}
	errch  chan error
	msgch  chan *pb.SendRequest
}

func (srv *server) Send(s pb.SenderService_SendServer) error {
	srv.log.Info("Got a send streaming request")
	ctx := s.Context()

	go srv.receive(s)
	go srv.send(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-srv.errch:
			return fmt.Errorf("failed to send msg: %w", err)
		case <-srv.donech:
			s.SendAndClose(&emptypb.Empty{})
			return nil
		}
	}
}

func (srv *server) receive(s pb.SenderService_SendServer) {
	for {
		msg, err := s.Recv()
		if errors.Is(err, io.EOF) {
			srv.donech <- struct{}{}
			return
		}

		if err != nil {
			srv.errch <- err
			return
		}

		srv.msgch <- msg
	}
}

func (srv *server) send(ctx context.Context) {
	for msg := range srv.msgch {
		srv.log.Debug("sending message...", "msg", msg.Message)
		if err := srv.sender.Send(ctx, sender.Message{Message: msg.Message}); err != nil {
			srv.errch <- fmt.Errorf("failed to send a message: %w", err)
			return
		}
	}
}

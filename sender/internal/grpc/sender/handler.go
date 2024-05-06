package sender

import (
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

func (srv *server) Send(ctx context.Context, s *pb.SendRequest) (*emptypb.Empty, error) {
	err := srv.sender.Send(ctx, sender.Message{Message: s.Message})
	return &emptypb.Empty{}, err
}

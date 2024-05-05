package translator

import (
	"context"
	"log/slog"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/translator"
	"google.golang.org/grpc"
)

func Register(srv *grpc.Server, t Translator, l *slog.Logger) {
	pb.RegisterTranslateServiceServer(srv, &server{
		translator: t,
		log:        l,
	})
}

type Translator interface {
	Translate(ctx context.Context, msg string) (string, error)
}

type server struct {
	pb.UnimplementedTranslateServiceServer
	translator Translator
	log        *slog.Logger
}

func (srv *server) Translate(
	ctx context.Context,
	req *pb.TranslateRequest,
) (*pb.TranslateResponse, error) {
	srv.log.Debug("incoming request")
	msg, err := srv.translator.Translate(ctx, req.Message)
	if err != nil {
		srv.log.Error("failed to transalte msg", "err", err)
		return nil, err
	}

	return &pb.TranslateResponse{Message: msg}, nil
}

package translator

import (
	"context"
	"errors"
	"fmt"
	"io"

	pb "github.com/hrvadl/goweekly/protos/gen/go/v1/translator"
	"google.golang.org/grpc"
)

func Register(srv *grpc.Server, t Translator) {
	pb.RegisterTranslateServiceServer(srv, &server{
		translator: t,
	})
}

type Translator interface {
	Translate(ctx context.Context, msg string) (string, error)
}

type server struct {
	pb.UnimplementedTranslateServiceServer
	translator Translator
}

func (srv *server) Translate(s pb.TranslateService_TranslateServer) error {
	ctx := s.Context()
	doneCh := make(chan struct{})
	errCh := make(chan error)

	go srv.handleStream(s, doneCh, errCh)

	select {
	case <-ctx.Done():
		return errors.New("timeout exceeded")
	case err := <-errCh:
		return fmt.Errorf("failed to translate articles: %w", err)
	case <-doneCh:
		return nil
	}
}

func (srv *server) handleStream(
	s pb.TranslateService_TranslateServer,
	doneCh chan<- struct{},
	errCh chan<- error,
) {
	ctx := s.Context()
	breakCh := make(chan struct{})

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
				msg, err := srv.translator.Translate(ctx, msg.Message)
				if err != nil {
					breakCh <- struct{}{}
					errCh <- err
				}
				if err := s.Send(&pb.TranslateRequest{Message: msg}); err != nil {
					breakCh <- struct{}{}
					errCh <- err
				}
			}()
		}
	}
}

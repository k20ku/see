package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// signal handler to do graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start the http request handle in another goroutine!
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				// ErrServerClosed is an expected error during a shutdown
				return nil
			}
			log.Printf("see server failed to close")
			return err
		}

		return nil
	})

	// Wait for an exit signal from the channel
	<-ctx.Done()

	shutdownCtx := context.Background()
	// and shutdown immediately after this server process got it!
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("see server failed to shutdown")
		return err
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	// wait for the another Goroutine to terminate which handles the incoming http request
	return eg.Wait()
}

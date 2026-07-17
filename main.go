package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k20ku/see/config"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	// OS signal handler
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// init config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("see server get config: %v", err)
	}

	// listen on port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("see server: failed to listen on port %d: %v", cfg.Port, err)
	}

	srv := http.Server{
		// 引数で受け取ったnet.Listenerを使うので
		// Addrフィールドは使用しない
		// Addr: server's port addr will be injected from the given net.Listener as args
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("see server: accepted request from %s.\n", r.URL.Path[1:])
			// コマンドラインでSIGINTを実験するため
			time.Sleep(5 * time.Second)
			if _, err := fmt.Fprintf(w,
				"Hello %s!\nYour User-Agent: %s\n", r.URL.Path[1:], r.Header.Get("User-Agent"),
			); err != nil {

				log.Printf("see server: failed to respond to %s.\n", r.URL.Path[1:])
				return
			}
			log.Printf("see server responds to %s.\n", r.URL.Path[1:])
		}),
	}

	// start the http request handle in another goroutine!
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// http.ErrServerClosed は異常系ではない
		// なぜならhttp.Server.Shutdownは正常に終了したことを示すので
		// 'Http.ErrServerClosed' is an expected error during a shutdown
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("see server failed to close")
			return err
		}

		return nil
	})

	log.Printf("see server: listen on port %d", cfg.Port)

	// Wait for an exit signal from the channel
	<-ctx.Done()

	shutdownCtx := context.Background()
	// and shutdown immediately after this server process got it!
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("see server failed to shutdown")
		return err
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	// wait for the another Goroutine to terminate which handles the incoming http request
	return eg.Wait()
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("see server failed to terminate: %s", err)
	}
}

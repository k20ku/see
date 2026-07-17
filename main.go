package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/k20ku/see/config"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("see server get config: %v", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("see server: failed to listen on port %d: %v", cfg.Port, err)
	}

	srv := http.Server{
		// 引数で受け取ったnet.Listenerを使うので
		// Addrフィールドは使用しない
		// Addr: server's port addr will be injected from the given net.Listener as args
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := fmt.Fprintf(w,
				"Hello %s!\nYour User-Agent: %s\n", r.URL.Path[1:], r.Header.Get("User-Agent"),
			); err != nil {

				log.Printf("see server: failed to respond to %s.\n", r.URL.Path[1:])
				return
			}
			log.Printf("see server: accepted request from %s.\n", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	// start the http request handle in another goroutine!
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

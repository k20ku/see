package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context, l net.Listener) error {
	srv := http.Server{
		// 引数で受け取ったnet.Listenerを使うので
		// Addrフィールドは使用しない
		// Addr: server's port addr will be injected from the given net.Listener as args
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
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
			log.Panicf("see server failed to close: %s", err)
			return err
		}

		return nil
	})

	// Wait for an exit signal from the channel
	<-ctx.Done()

	shutdownCtx := context.Background()
	// and shutdown immediately after this server process got it!
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Panicf("see server failed to shutdown: %+v", err)
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	// wait for the another Goroutine to terminate which handles the incoming http request
	return eg.Wait()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("see sever need port number\n")
	}

	// init the Listener
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("see server failed to listen port %s: %+v", p, err)
	}

	ctx := context.Background()
	if err := run(ctx, l); err != nil {
		log.Fatalf("see server failed to terminate: %s", err)
	}
}

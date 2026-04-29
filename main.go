package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	srv := http.Server{
		Addr: ":18080",
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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Printf("see server failed to terminate: %s", err)
	}
}

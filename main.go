package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/k20ku/see/config"
	"github.com/k20ku/see/server"
)

func run(ctx context.Context) error {
	// init config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("see server get config: %w", err)
	}

	// listen on port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("see server: failed to listen on port %d: %w", cfg.Port, err)
	}

	mux := server.NewMux()
	srv := server.NewServer(l, mux)

	// start the http request handle in another goroutine!
	log.Printf("see server: listen on port %d", cfg.Port)
	return srv.Run(ctx)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("see server failed to terminate: %s", err)
	}
}

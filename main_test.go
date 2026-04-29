package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	// init the TCP Listener for See Server
	// ポート番号0を指定すると，ポート番号は自動的に設定される
	// ポート番号が固定されていると，他のAppがそのポートを使ってる場合競合が起きるため
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port: %v", err)
	}
	// run the server with a cancel context in another proess
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	// this test client send http request
	in := "World"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// どんなポートでlistenしているのかログ
	t.Logf("try request to %q", url)
	// wait enough time for the server to start
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to get: %+v", err)
	}
	defer resp.Body.Close()

	// get the response from the server which started at step.1
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %+v", err)
	}

	// validate the gotten http response
	want := fmt.Sprintf("Hello %s!", in)
	if string(got) != want {
		t.Errorf("unexpected response, want=%q, got=%q", want, got)
	}

	// send a cancel signal to server in background
	cancel()
	// Did server successfully shutdown?
	if err := eg.Wait(); err != nil {
		t.Fatalf("server failed to succesefully shutdown: %+v", err)
	}
}

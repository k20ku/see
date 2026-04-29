package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	// 1. run the server with a cancel context in another proess
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	// 2. this test client send http request
	in := "World"
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Fatalf("failed to get: %+v", err)
	}
	defer resp.Body.Close()

	// 3. get the response from the server which started at step.1
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %+v", err)
	}

	// 4. validate the gotten http response
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

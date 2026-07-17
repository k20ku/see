package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestServerRun(t *testing.T) {
	// init the TCP Listener for See Server
	// ポート番号0を指定すると，ポート番号は自動的に設定される
	// ポート番号が固定されていると，他のAppがそのポートを使ってる場合競合が起きるため
	l, err := net.Listen("tcp", "localhost:0")
	require.NoErrorf(t, err, "failed to listen on port %d", 0)

	// run the server with a cancel context in another process
	ctx, cancel := context.WithCancel(context.Background())

	// init server
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:]); err != nil {
			t.Logf("failed to respond to %s: %v", r.URL.Path, err)
		}
	})

	// run server
	srv := NewServer(l, mux)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return srv.Run(ctx)
	})

	// this test client send http GET request
	in := "World"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	client := http.Client{Timeout: 60 * time.Second}

	t.Logf("try request to %q", url)
	// wait enough time for the server to start
	time.Sleep(100 * time.Millisecond)

	// request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err, "prepare GET request")
	rawReq, err := httputil.DumpRequest(req, true)
	require.NoError(t, err, "dump request: %s", rawReq)

	resp, err := client.Do(req)
	require.NoError(t, err, "client sends request")
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Log("resp.Body failed to close")
		}
	}()

	respRaw, err := httputil.DumpResponse(resp, true)
	require.NoErrorf(t, err, "client recieved the response: %s", respRaw)
	// get the response from the server which started at step.1
	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "read response body")

	// validate the gotten http response
	want := fmt.Sprintf("Hello %s!", in)
	require.Equal(t, string(got), want, "validate gotten response")

	// send a cancel signal to server in background
	cancel()

	err = eg.Wait()
	require.NoError(t, err, "shutdown server")
}

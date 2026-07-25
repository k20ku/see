# see

## how to use

```bash
make build
make up
```

Open another terminal, run below.

```bash
curl -i -XGET localhost:18080/health
```

### Response

```http
$ curl -i -w"\n" -XGET localhost:18080/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 24 Jul 2026 14:19:41 GMT
Content-Length: 17

{"status" : "ok"}
```

## Futures

- Graceful Shutdown

    ```bash
    docker run -p 28080:80 k20ku/see:latest
    ```

    Even if we send `SIGINT` immediately after server have received request, server exits after sending the response.

    ```log
    2026/07/17 13:50:25 see server: listen on port 80
    2026/07/17 13:50:26 see server: accepted request from hello.
    ^C2026/07/17 13:50:31 see server responds to hello.
    ```

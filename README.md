# see

## how to use

```bash
make build
docker run -p 18080:80 k20ku/see:latest
```

Open another terminal, run below.

```bash
curl -i http://localhost:18080/World
```

### Response

```http
HTTP/1.1 200 OK
Date: Fri, 08 May 2026 18:59:38 GMT
Content-Length: 41
Content-Type: text/plain; charset=utf-8

Hello World!
Your User-Agent: curl/8.5.0
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

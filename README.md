# see

## how to use

```bash
go run . localhost 18080
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

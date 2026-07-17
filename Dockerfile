# the container that creates the binary included in the deployed container
FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
        -trimpath \
        -ldflags="-s -w" \
        -o app .

# 
FROM golang:1.26-bookworm AS dev
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]

# 
FROM gcr.io/distroless/static-debian12 AS deploy
COPY --from=builder /app/app .
USER nonroot:nonroot
ENTRYPOINT ["./app"]


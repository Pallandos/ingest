FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod .
COPY cmd/ cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/ingest-server ./cmd/server

# ---

FROM scratch

COPY --from=builder /bin/ingest-server /ingest-server

EXPOSE 8080
ENV LISTEN_ADDR=:8080

ENTRYPOINT ["/ingest-server"]
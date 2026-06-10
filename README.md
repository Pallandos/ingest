# ingest

Tiny data collector server, written in Go 

## Usage

`ingest` is a server exposing `POST` endpoint to colelct data. 

For example : 

```sh
curl -s -X POST http://localhost:8080/logs \
  -H 'Content-Type: application/json' \
  -d '{"level":"info","message":"hello","service":"myapp"}'
```

## Format 

TODO
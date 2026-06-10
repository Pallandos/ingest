# ingest

Tiny data collector server, written in Go 

## Usage

`ingest` is a server exposing endpoint to collect data. You can send data to a **channel** which is like a subject. Each channel has its own data directory and endpoint. 

To post data to a channel reach the following endpoint : 

    http://myserver/data/{channel}

From now you can pass your payload and data will be recorded inside the `/data` docker volume. 

## Example :

```sh
curl -X POST http://localhost:8080/data/temparature \
-H "Content-Type: application/json" \
-d '{"value": 22.5, "unit": "Celsius"}'
```

This will add the following line to `/data/temparature` :

```json
{"received_at":"2026-06-10T20:25:07.629794925Z","remote_addr":"172.18.0.1:43004","payload":{"value":22.5,"unit":"Celsius"}}
```
# ingest

Tiny data collector server, written in Go 

## Installation 

You can use either use `docker run` or `docker compose`. See bellow for examples. 
<details>
<summary> `docker compose` </summary>

Create a `docker-compose.yml` file, for example :

```yml
services:
  ingest-server:
    image: ghcr.io/pallandos/ingest:latest
    ports:
      - "8080:8080"
    environment:
      LISTEN_ADDR: ":8080"
      DATA_DIR: "/data"
    volumes:
      - ./data:/data
    restart: unless-stopped
    logging:
      driver: json-file
      options:
        max-size: "50m"
        max-file: "5"
```

And then run 

    docker compose up -d

</details>

<details>
<summary> `docker run` </summary>

```sh
docker run -d \
  --name ingest-server \
  -p 8080:8080 \
  -e LISTEN_ADDR=":8080" \
  -e DATA_DIR="/data" \
  -v "$(pwd)/data:/data" \
  --restart unless-stopped \
  --log-driver json-file \
  --log-opt max-size=50m \
  --log-opt max-file=5 \
  ghcr.io/pallandos/ingest:latest
```

</details>



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
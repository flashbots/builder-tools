# builder-tools

WIP Toolbox

- [Create ECDSA keypair](cmd/ecdsa-gen/main.go)
- [Create TLS certificate + key (PEM format)](cmd/tls-gen/main.go)
- [Server using custom TLS certificate](cmd/https-server/main.go)
- [Client allowing only server using the custom TLS certificate](cmd/https-client/main.go)
- [Status API server, with ability for recording and querying events](cmd/status-api/)

Canonical way to collapse a TDX [measurements.json](docs/measurements.json) file into a single hash, in a reproducible way:

```bash
cat measurements.json | jq --sort-keys --compact-output --join-output | sha256sum
```

## Usage

```bash
# create the TLS cert (cert.pem) and key (key.pem)
$ go run cmd/tls-gen/main.go --host 127.0.0.1,localhost

# run the server (serving the created TLS cert)
$ go run cmd/https-server/main.go

# check with curl
$ curl --cacert cert.pem https://127.0.0.1:8080

# run the client (allowing only server with that specific TLS cert)
$ go run cmd/https-client/main.go
```

### Status API server

The status api server is used to record and query events. Events can be added through local named pipe (file `pipe.fifo`), or through HTTP API.

```bash
# Start the server
$ go run cmd/status-api/*

# Add events
$ echo "hello world" > pipe.fifo
$ curl localhost:8082/api/v1/new_event?message=this+is+a+test

# Query events (timestamp in UTC)
$ curl -s localhost:8082/api/v1/events | jq -r  '(.[] | [.received_at, .message]) | @tsv'
2024-09-24T10:45:50.774339Z     hello world
2024-09-24T10:46:02.01221Z      this is a test
```

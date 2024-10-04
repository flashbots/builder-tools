# builder-tools

WIP Toolbox

- [Create ECDSA keypair](cmd/ecdsa-gen/main.go)
- [Create TLS certificate + key (PEM format)](cmd/tls-gen/main.go)
- [Server using custom TLS certificate](cmd/https-server/main.go)
- [Client allowing only server using the custom TLS certificate](cmd/https-client/main.go)
- [Status API server, with ability for recording and querying events](cmd/status-api/)

---

## Single-Hash TDX Measurement

We use this method to collapse a TDX [measurements JSON object](docs/measurements.json) into a single SHA256 hash, in a reproducible way (sort keys, remove whitespace, print without newline):

```bash
cat docs/measurements.json | jq --sort-keys --compact-output --join-output | sha256sum
59952c2557da91cdafb9361d41d6971fc2b0be1a85ad57134e38714b266ff581  -
```

---

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

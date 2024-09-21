# builder-tools

WIP Toolbox

- [Create ECDSA keypair](cmd/ecdsa-gen/main.go)
- [Create TLS certificate + key (PEM format)](cmd/tls-gen/main.go)
- [Server using custom TLS certificate](cmd/https-server/main.go)
- [Client allowing only server using the custom TLS certificate](cmd/https-client/main.go)
- [Status API server for recording and querying events](cmd/status-api/)

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

Status api server that can be used to record and query events.

```bash
# Start the server
$ go run cmd/status-api/*

# Add events
$ curl localhost:8082/api/v1/new_event?message=111
$ curl localhost:8082/api/v1/new_event?message=222

# Query events
$ curl -s localhost:8082/api/v1/events | jq
[
  {
    "ReceivedAt": "2024-09-21T07:51:47.709512Z",
    "Message": "111"
  },
  {
    "ReceivedAt": "2024-09-21T07:51:49.294175Z",
    "Message": "222"
  }
]
```


## Next Steps

These partly overlap with https://github.com/flashbots/cvm-reverse-proxy:
- Server that verifies client-side aTLS certificate
- Client that sends client-side aTLS certificate
- One server that exposes an aTLS endpoint to serve the local TLS cert, and another server that exposes a TLS endpoint

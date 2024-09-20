# builder-tools

WIP Toolbox

- [Create ECDSA keypair](cmd/ecdsa-gen/main.go)
- [Create TLS certificate + key (PEM format)](cmd/tls-gen/main.go)
- [Server using custom TLS certificate](cmd/https-server/main.go)
- [Client allowing only server using the custom TLS certificate](cmd/https-client/main.go)

## Usage

```bash
# create the TLS cert
$ go run cmd/tls-gen/main.go --host foo.com,0.0.0.0,127.0.0.1

# run the server (serving the created TLS cert)
$ go run cmd/https-server/main.go

# check with curl
$ curl --insecure https://localhost:8080

# run the client (allowing only server with that specific TLS cert)
$ go run cmd/https-client/main.go
```

## Next Steps

These partly overlap with https://github.com/flashbots/cvm-reverse-proxy:
- Server that verifies client-side aTLS certificate
- Client that sends client-side aTLS certificate
- One server that exposes an aTLS endpoint to serve the local TLS cert, and another server that exposes a TLS endpoint

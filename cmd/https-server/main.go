package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/flashbots/go-template/common"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

var flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:  "listen-addr",
		Value: "0.0.0.0:8080",
		Usage: "address to serve certificate on",
	},
}

func main() {
	app := &cli.App{
		Name:   "https-server",
		Usage:  "Server with custom cert",
		Flags:  flags,
		Action: runCli,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runCli(cCtx *cli.Context) error {
	listenAddr := cCtx.String("listen-addr")
	certFile := "cert.pem"
	keyFile := "key.pem"

	log := common.SetupLogger(&common.LoggingOpts{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	log.Info("Starting HTTPS server", "addr", listenAddr)
	if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Error("proxy exited", "err", err)
		return err
	}
	return nil
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "checkcheck\n")
}

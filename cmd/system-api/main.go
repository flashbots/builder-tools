package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flashbots/system-api/common"
	cli "github.com/urfave/cli/v2" // imports as package "cli"
)

var flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:  "listen-addr",
		Value: "0.0.0.0:8082",
		Usage: "address to serve certificate on",
	},
	&cli.StringFlag{
		Name:  "pipe-file",
		Value: "pipe.fifo",
		Usage: "filename for named pipe (for sending events into this service)",
	},
}

func main() {
	app := &cli.App{
		Name:    "system-api",
		Usage:   "HTTP API for status events",
		Version: common.Version,
		Flags:   flags,
		Action:  runCli,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runCli(cCtx *cli.Context) error {
	listenAddr := cCtx.String("listen-addr")
	pipeFile := cCtx.String("pipe-file")

	log := common.SetupLogger(&common.LoggingOpts{
		Version: common.Version,
	})

	// Setup and start the server (in the background)
	server, err := NewServer(&HTTPServerConfig{
		ListenAddr:   listenAddr,
		Log:          log,
		PipeFilename: pipeFile,
	})
	if err != nil {
		return err
	}
	go server.Start()

	// Wait for signal, then graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Error("HTTP shutdown error", "err", err)
		return err
	}
	return nil
}

package main

import (
	"log"
	"os"

	"github.com/flashbots/go-template/common"
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
		Name:   "status-logger-api",
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
	pipeFile := cCtx.String("pipe-file")

	log := common.SetupLogger(&common.LoggingOpts{})
	server, err := NewServer(&HTTPServerConfig{
		ListenAddr:   listenAddr,
		Log:          log,
		PipeFilename: pipeFile,
	})
	if err != nil {
		return err
	}
	server.Start()
	return nil
}

package main

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/flashbots/go-template/common"
	cli "github.com/urfave/cli/v2" // imports as package "cli"
)

// var flags []cli.Flag = []cli.Flag{
// 	&cli.BoolFlag{
// 		Name:  "log-json",
// 		Value: false,
// 		Usage: "log in JSON format",
// 	},
// 	&cli.BoolFlag{
// 		Name:  "log-debug",
// 		Value: false,
// 		Usage: "log debug messages",
// 	},
// }

func main() {
	app := &cli.App{
		Name:  "ecdsa-gen",
		Usage: "Create ECDSA keypair",
		// Flags:  flags,
		Action: runCli,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runCli(cCtx *cli.Context) error {
	log := common.SetupLogger(&common.LoggingOpts{})

	ecdsaPrivkeyHex, ecdsaAddressHex, err := genECDSA()
	if err != nil {
		return err
	}

	log.Info("ECDSA keypair generated", "private", ecdsaPrivkeyHex, "address", ecdsaAddressHex)
	return nil
}

func genECDSA() (privkey, address string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privkeyHex := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	addressHex := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return privkeyHex, addressHex, nil
}

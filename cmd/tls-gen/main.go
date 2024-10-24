// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

// Source: https://go.dev/src/crypto/tls/generate_cert.go
// See also: https://gist.github.com/denji/12b3a568f092ab951456

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/flashbots/system-api/crypto"
)

var (
	host       = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	validFrom  = flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	validFor   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
	isCA       = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	rsaBits    = flag.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	ecdsaCurve = flag.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
	ed25519Key = flag.Bool("ed25519", false, "Generate an Ed25519 key")
)

func main() {
	flag.Parse()

	if len(*host) == 0 {
		log.Fatalf("Missing required --host parameter")
	}

	opts := crypto.GenTLSOpts{
		Organisation: "Flashbots",
		Host:         *host,
		ValidFrom:    *validFrom,
		ValidFor:     *validFor,
		IsCA:         *isCA,
		RsaBits:      *rsaBits,
		EcdsaCurve:   *ecdsaCurve,
		Ed25519Key:   *ed25519Key,
	}
	cert, priv, err := crypto.GenTLS(&opts)
	if err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	// Write to file: cert.pem
	certOut, err := os.Create("cert.pem")
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	if _, err := certOut.Write([]byte(cert)); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing cert.pem: %v", err)
	}
	log.Print("wrote cert.pem\n")

	// Write to file: key.pem
	keyOut, err := os.Create("key.pem")
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
	}
	if _, err := keyOut.Write([]byte(priv)); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	log.Print("wrote key.pem\n")
}

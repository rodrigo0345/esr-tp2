package config

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/rodrigo0345/esr-tp2/server/certs"
)

func GenerateTLS() *tls.Config {

	certFile, keyFile, err := certs.GenerateTLSFiles()

	if err != nil {
		log.Fatalf("Failed to generate TLS files: %v", err)
	}
	defer os.Remove(certFile)
	defer os.Remove(keyFile)

	// Load the TLS configuration
	tlsConfig, err := certs.GenerateTLSConfig(certFile, keyFile)

	if err != nil {
		log.Fatalf("Failed to load TLS configuration: %v", err)
	}
	return tlsConfig
}

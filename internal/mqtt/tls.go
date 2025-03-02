package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func NewTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	var certs []tls.Certificate
	if certFile != "" && keyFile != "" {
		clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate and key: %w", err)
		}
		certs = append(certs, clientCert)
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: certs,
	}

	return tlsConfig, nil
}

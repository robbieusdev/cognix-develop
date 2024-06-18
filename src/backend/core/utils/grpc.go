package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type (
	// TLSConfig is a type that represents
	// common TLS configuration
	TLSConfig struct {
		CACertificate   string
		CertificateFile string
		KeyFile         string
	}
	// credentialsConstructor is a type
	// that represents credentials creator function signature
	credentialsConstructor func(
		certPool *x509.CertPool, certificate *tls.Certificate,
	) (credentials.TransportCredentials, error)
	// certificatePoolConstructor is a type
	// that represents certificate pool creator function signature
	certificatePoolConstructor func(
		caCertFile, certFile, keyFile []byte,
	) (*x509.CertPool, *tls.Certificate, error)
)

// NewServerTLSCredentials creates a new instance of the TransportCredentials.
func NewServerTLSCredentials(cfg *TLSConfig) (credentials.TransportCredentials, error) {
	return withCertificatePool(cfg,
		newTLSCertificateWithPool,
		func(certPool *x509.CertPool, certificate *tls.Certificate) (credentials.TransportCredentials, error) {
			return credentials.NewTLS(&tls.Config{
				MinVersion:   tls.VersionTLS12,
				ClientAuth:   tls.RequireAndVerifyClientCert,
				Certificates: []tls.Certificate{*certificate},
				ClientCAs:    certPool,
			}), nil
		})
}

// NewClientTLSCredentials creates a new instance of the TransportCredentials.
func NewClientTLSCredentials(cfg *TLSConfig) (credentials.TransportCredentials, error) {
	return withCertificatePool(cfg,
		newTLSCertificateWithPool,
		func(certPool *x509.CertPool, certificate *tls.Certificate) (credentials.TransportCredentials, error) {
			return credentials.NewTLS(&tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{*certificate},
				RootCAs:      certPool,
			}), nil
		})
}

// withCertificatePool decorates constructor with certificates preloading
func withCertificatePool(
	cfg *TLSConfig,
	constructCertificate certificatePoolConstructor,
	constructCredentials credentialsConstructor,
) (credentials.TransportCredentials, error) {
	caRootCert, err := ioutil.ReadFile(cfg.CACertificate)
	if err != nil {
		return nil, err
	}
	certPEMBlock, err := ioutil.ReadFile(cfg.CACertificate)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := ioutil.ReadFile(cfg.CACertificate)
	if err != nil {
		return nil, err
	}
	certPool, certificate, err := constructCertificate(caRootCert, certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	return constructCredentials(certPool, certificate)
}

// newTLSCertificateWithPool creates a new pair of CertPool and Certificate.
func newTLSCertificateWithPool(caCertFile, certFile, keyFile []byte) (*x509.CertPool, *tls.Certificate, error) {
	certificate, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(caCertFile)
	if !ok {
		return nil, nil, fmt.Errorf("failed to append certs")
	}
	return certPool, &certificate, nil
}

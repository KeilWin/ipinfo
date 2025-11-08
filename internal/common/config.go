package common

import (
	"crypto/tls"
	"fmt"
)

type Config interface {
	Load() error
	Check() error
	NewVariableName(name string) string
}

func NewBasePrefix(appPrefix string, componentPrefix string) string {
	return fmt.Sprintf("%s_%s", appPrefix, componentPrefix)
}

func NewCipherSuites() []uint16 {
	return []uint16{
		// Safe ciphers from https://www.ssllabs.com/
		// TLS 1.2
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		// TLS 1.3
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
	}
}

func NewCurvePreferences() []tls.CurveID {
	return []tls.CurveID{
		tls.CurveP521,
		tls.X25519,
	}
}

func NewNextProtos() []string {
	return []string{"h2", "http/1.1"}
}

func NewTlsConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,

		CipherSuites:     NewCipherSuites(),
		CurvePreferences: NewCurvePreferences(),

		SessionTicketsDisabled: false,
		Renegotiation:          tls.RenegotiateNever,

		InsecureSkipVerify: true,

		NextProtos: NewNextProtos(),
	}
}

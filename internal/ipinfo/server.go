package app

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/utils"
)

type ProtocolType string

const (
	ProtocolHTTP  ProtocolType = "http"
	ProtocolHTTPS ProtocolType = "https"
)

const componentName = "SERVER"

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

		NextProtos: NewNextProtos(),
	}
}

func NewAppServer(handler *http.ServeMux, appCfg *ServerConfig) *http.Server {
	return &http.Server{
		Addr:              appCfg.Addr,
		Handler:           handler,
		TLSConfig:         NewTlsConfig(),
		MaxHeaderBytes:    appCfg.MaxHeaderBytes,
		ReadTimeout:       appCfg.ReadTimeout,
		ReadHeaderTimeout: appCfg.ReadHeaderTimeout,
		WriteTimeout:      appCfg.WriteTimeout,
		IdleTimeout:       appCfg.IdleTimeout,
	}
}

func CheckLoadConfigError(err error, name string) bool {
	return utils.CheckLoadConfigError(err, name, componentName)
}

type ServerConfig struct {
	common.Config
	BasePrefix string

	Addr              string
	Protocol          ProtocolType
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	CertFile string
	KeyFile  string
}

func (p *ServerConfig) NewVariableName(name string) string {
	return fmt.Sprintf("%s_%s", p.BasePrefix, name)
}

func (p *ServerConfig) Load() error {
	var err error
	var hasError bool

	addrName := p.NewVariableName("ADDR")
	p.Addr = os.Getenv(addrName)
	protocolName := p.NewVariableName("PROTOCOL")
	p.Protocol = ProtocolType(os.Getenv(protocolName))
	maxHeaderBytesName := p.NewVariableName("MAX_HEADER_BYTES")
	p.MaxHeaderBytes, err = strconv.Atoi(os.Getenv(maxHeaderBytesName))
	hasError = CheckLoadConfigError(err, maxHeaderBytesName)

	readTimeoutName := p.NewVariableName("READ_TIMEOUT")
	readTimeout, err := strconv.Atoi(os.Getenv(readTimeoutName))
	hasError = CheckLoadConfigError(err, readTimeoutName)
	p.ReadTimeout = time.Duration(readTimeout) * time.Second

	readHeaderTimeoutName := p.NewVariableName("READ_HEADER_TIMEOUT")
	readHeaderTimeout, err := strconv.Atoi(os.Getenv(readHeaderTimeoutName))
	hasError = CheckLoadConfigError(err, readHeaderTimeoutName)
	p.ReadHeaderTimeout = time.Duration(readHeaderTimeout) * time.Second

	writeTimeoutName := p.NewVariableName("WRITE_TIMEOUT")
	writeTimeout, err := strconv.Atoi(os.Getenv(writeTimeoutName))
	hasError = CheckLoadConfigError(err, writeTimeoutName)
	p.WriteTimeout = time.Duration(writeTimeout) * time.Second

	idleTimeoutName := p.NewVariableName("IDLE_TIMEOUT")
	idleTimeout, err := strconv.Atoi(os.Getenv(idleTimeoutName))
	hasError = CheckLoadConfigError(err, idleTimeoutName)
	p.IdleTimeout = time.Duration(idleTimeout) * time.Second

	certFileName := p.NewVariableName("CERT_FILE")
	p.CertFile = os.Getenv(certFileName)
	keyFileName := p.NewVariableName("KEY_FILE")
	p.KeyFile = os.Getenv(keyFileName)

	if hasError {
		return errors.New("loading server config")
	}
	return nil
}

func (p *ServerConfig) Check() error {
	return nil
}

func NewServerConfig(appPrefix string) *ServerConfig {
	return &ServerConfig{
		BasePrefix: common.NewBasePrefix(appPrefix, componentName),
	}
}

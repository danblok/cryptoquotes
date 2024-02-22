package config

import (
	"crypto/tls"
	"fmt"
)

// HTTPConfig is a config for HTTP API server.
type HTTPConfig struct {
	host      string
	port      uint16
	tlsConfig *tls.Config
}

// NewHTTPConfig constructs a new Config for HTTP API server.
func NewHTTPConfig(host string, port uint16, tlsConfig *tls.Config) *HTTPConfig {
	return &HTTPConfig{
		host:      host,
		port:      port,
		tlsConfig: tlsConfig,
	}
}

// URL returns the full URL of HTTP API server.
func (c *HTTPConfig) URL() string {
	if c.tlsConfig == nil {
		return fmt.Sprintf("http://%s:%d", c.host, c.port)
	}
	return fmt.Sprintf("https://%s:%d", c.host, c.port)
}

// TLSConfig returns the certificate of the HTTP API server.
func (c *HTTPConfig) TLSConfig() *tls.Config {
	return c.tlsConfig
}

// Port returns the port to which the HTTP API server listens to.
func (c *HTTPConfig) Port() uint16 {
	return c.port
}

// Host returns the host of the HTTP API server.
func (c *HTTPConfig) Host() string {
	return c.host
}

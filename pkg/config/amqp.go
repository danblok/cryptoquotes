package config

import (
	"crypto/tls"
	"fmt"
)

// AMQPConfig is a config for a message broker.
type AMQPConfig struct {
	login    string
	password string
	host     string
	port     uint16
	cert     *tls.Certificate
}

// NewAMQPConfig constructs a new Config for a message broker.
func NewAMQPConfig(login, password, host string, port uint16, cert *tls.Certificate) *AMQPConfig {
	return &AMQPConfig{
		login:    login,
		password: password,
		host:     host,
		port:     port,
		cert:     cert,
	}
}

// URL returns the full URL of the message broker server.
func (c *AMQPConfig) URL() string {
	if c.cert == nil {
		return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.login, c.password, c.host, c.port)
	}
	return fmt.Sprintf("amqps://%s:%s@%s:%d/", c.login, c.password, c.host, c.port)
}

// Cert returns the certificate of the message broker server.
func (c *AMQPConfig) Cert() *tls.Certificate {
	return c.cert
}

// Port returns the port to which the message broker server listens to.
func (c *AMQPConfig) Port() uint16 {
	return c.port
}

// Host returns the host of the message broker server.
func (c *AMQPConfig) Host() string {
	return c.host
}

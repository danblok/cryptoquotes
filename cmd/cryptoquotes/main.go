package main

import (
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"github.com/danblok/cryptoquotes/internal/api"
	"github.com/danblok/cryptoquotes/internal/logging"
	"github.com/danblok/cryptoquotes/internal/messages"
	"github.com/danblok/cryptoquotes/pkg/config"
	"github.com/danblok/cryptoquotes/pkg/source"
)

func main() {
	m := messages.New(config.NewAMQPConfig(os.Getenv("MQ_LOGIN"), os.Getenv("MQ_PASSWORD"), "mq", 5672, nil))
	ch, err := m.EmittingChannel()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	svc := source.NewCoinMarketCapSource(os.Getenv("COINMARKETCAP_API_KEY"))
	svc = logging.New(svc, ch)
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Fatal("couldn't parsed the port")
	}
	serverCert, err := tls.LoadX509KeyPair("/run/secrets/api_cert", "/run/secrets/api_key")
	if err != nil {
		log.Fatal(err)
	}

	tlsServerConf := &tls.Config{Certificates: []tls.Certificate{serverCert}}

	srv := api.NewHTTPServer(svc, config.NewHTTPConfig("localhost", uint16(port), tlsServerConf))
	if err := srv.ListenAndServeTLS(); err != nil {
		log.Fatal(err)
	}
}

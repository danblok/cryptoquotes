package source

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danblok/cryptoquotes/pkg/types"
)

// UnmarshalFunc parses a response body and extracts quotes from the parsed body.
// The client should call Close method themselves.
type UnmarshalFunc func(io.ReadCloser) ([]types.Quote, error)

// HTTPSource allows to fetch crypto quotes from 3rd party API.
type HTTPSource struct {
	client    *http.Client
	url       string
	headers   map[string]string
	vals      url.Values
	unmarshal UnmarshalFunc
}

// NewHTTPSource contsructs a new HTTPSource instance.
func NewHTTPSource(
	url string,
	headers map[string]string,
	vals url.Values,
	unmarshal UnmarshalFunc,
) types.QuotesSource {
	return &HTTPSource{
		url:       url,
		headers:   headers,
		vals:      vals,
		unmarshal: unmarshal,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// NewHTTPSourceTLS contsructs a new HTTPSource instance with TLS support.
func NewHTTPSourceTLS(
	url string,
	headers map[string]string,
	vals url.Values,
	cert []byte,
	unmarshal UnmarshalFunc,
) (types.QuotesSource, error) {
	// load certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("couldn't append cert from file")
	}

	tlsCnf := &tls.Config{
		RootCAs: certPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsCnf,
	}
	return &HTTPSource{
		url:       url,
		headers:   headers,
		vals:      vals,
		unmarshal: unmarshal,
		client: &http.Client{
			Timeout:   3 * time.Second,
			Transport: transport,
		},
	}, nil
}

// Quotes fetches crypto quotes.
func (s *HTTPSource) Quotes(ctx context.Context) ([]types.Quote, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.url, nil)
	if err != nil {
		return nil, err
	}
	for hkey, h := range s.headers {
		req.Header.Add(hkey, h)
	}
	req.URL.RawQuery = s.vals.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	quotes, err := s.unmarshal(resp.Body)
	if err != nil {
		return nil, err
	}

	return quotes, nil
}

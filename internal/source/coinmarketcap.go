package source

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/danblok/cryptoquotes/pkg/types"
)

// NewCoinMarketCapSource returns HTTPSource of CoinMarketCap.
func NewCoinMarketCapSource(apiKey string, vals url.Values) types.QuotesSource {
	return NewHTTPSource(
		"https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest",
		map[string]string{
			"Accepts":           "application/json",
			"X-CMC_PRO_API_KEY": apiKey,
		},
		vals,
		func(rc io.ReadCloser) ([]types.Quote, error) {
			defer rc.Close()

			var body struct {
				Status any `json:"status"`
				Data   map[string][]struct {
					Name        string  `json:"name"`
					Symbol      string  `json:"symbol"`
					TotalSupply float64 `json:"total_supply"`
					Quote       struct {
						USD struct {
							Price float64 `json:"price"`
						} `json:"USD"`
					} `json:"quote"`
				} `json:"data"`
			}

			if err := json.NewDecoder(rc).Decode(&body); err != nil {
				return nil, err
			}

			var quotes []types.Quote
			for _, quote := range body.Data {
				quotes = append(quotes, types.Quote{
					Name:        quote[0].Name,
					Symbol:      quote[0].Symbol,
					TotalSupply: quote[0].TotalSupply,
					USD:         quote[0].Quote.USD.Price,
				})
			}

			return quotes, nil
		},
	)
}

// NewCoinMarketCapSourceTLS returns HTTPSource of CoinMarketCap with TLS enabled.
func NewCoinMarketCapSourceTLS(apiKey string, vals url.Values, cert []byte) (types.QuotesSource, error) {
	return NewHTTPSourceTLS(
		"https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest",
		map[string]string{
			"Accepts":           "application/json",
			"X-CMC_PRO_API_KEY": apiKey,
		},
		vals,
		cert,
		func(rc io.ReadCloser) ([]types.Quote, error) {
			defer rc.Close()

			var body struct {
				Status any `json:"status"`
				Data   map[string][]struct {
					Name        string  `json:"name"`
					Symbol      string  `json:"symbol"`
					TotalSupply float64 `json:"total_supply"`
					Quote       struct {
						USD struct {
							Price float64 `json:"price"`
						} `json:"USD"`
					} `json:"quote"`
				} `json:"data"`
			}

			if err := json.NewDecoder(rc).Decode(&body); err != nil {
				return nil, err
			}

			var quotes []types.Quote
			for _, quote := range body.Data {
				quotes = append(quotes, types.Quote{
					Name:        quote[0].Name,
					Symbol:      quote[0].Symbol,
					TotalSupply: quote[0].TotalSupply,
					USD:         quote[0].Quote.USD.Price,
				})
			}

			return quotes, nil
		},
	)
}

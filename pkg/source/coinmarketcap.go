package source

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danblok/cryptoquotes/pkg/types"
)

const coinMarketCapURL = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"

type coinMarketCapSource struct {
	apiKey string
}

type coinMarketCapBody struct {
	Status struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"status"`
	Data map[string][]struct {
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

// NewCoinMarketCapSource returns a new NewCoinMarketCapSource as QuotesSource.
func NewCoinMarketCapSource(apiKey string) types.QuotesSource {
	return &coinMarketCapSource{
		apiKey: apiKey,
	}
}

func (s *coinMarketCapSource) Quotes(ctx context.Context, names ...string) ([]types.Quote, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", coinMarketCapURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", s.apiKey)

	query := req.URL.Query()
	query.Add("symbol", strings.Join(names, ","))
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body coinMarketCapBody
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	var quotes []types.Quote
	for _, quote := range body.Data {
		if len(quote) > 0 {
			quotes = append(quotes, types.Quote{
				Name:        quote[0].Name,
				Symbol:      quote[0].Symbol,
				TotalSupply: quote[0].TotalSupply,
				USD:         quote[0].Quote.USD.Price,
			})
		}
	}

	return quotes, nil
}

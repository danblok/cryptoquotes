package types

import "context"

// QuotesSource returns Quotes.
type QuotesSource interface {
	Quotes(context.Context) ([]Quote, error)
}

// Quote contains info about Quote.
type Quote struct {
	Name        string
	Symbol      string
	USD         float64
	TotalSupply float64
}

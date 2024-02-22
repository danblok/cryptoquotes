package logging

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	mq "github.com/rabbitmq/amqp091-go"

	"github.com/danblok/cryptoquotes/pkg/types"
)

type loggingService struct {
	svc types.QuotesSource
	ch  *mq.Channel
}

// Log represents a logging record.
type Log struct {
	RequestID     string
	StartTime     time.Time
	RequestTime   time.Duration
	TransportType string
	RemoteAddr    string
	Error         string
	Quotes        []types.Quote
}

// New returns a new logging service.
func New(svc types.QuotesSource, ch *mq.Channel) types.QuotesSource {
	return &loggingService{
		svc: svc,
		ch:  ch,
	}
}

func (s *loggingService) Quotes(ctx context.Context, names ...string) (quotes []types.Quote, err error) {
	defer func(t time.Time) {
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		var buf bytes.Buffer
		err = gob.NewEncoder(&buf).Encode(Log{
			RequestID:     ctx.Value(types.RequestIDKey).(string),
			StartTime:     t,
			RequestTime:   time.Since(t),
			TransportType: ctx.Value(types.TransportTypeKey).(string),
			RemoteAddr:    ctx.Value(types.RemoteAddrKey).(string),
			Error:         errMsg,
			Quotes:        quotes,
		})
		if err != nil {
			return
		}

		err = s.ch.PublishWithContext(
			ctx,
			"logs",
			"",
			false,
			false,
			mq.Publishing{ContentType: "application/octet-stream", Body: buf.Bytes()},
		)
	}(time.Now())

	return s.svc.Quotes(ctx, names...)
}

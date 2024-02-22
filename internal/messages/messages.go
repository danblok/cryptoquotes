package messages

import (
	"bytes"
	"context"
	"encoding/gob"
	"log/slog"

	mq "github.com/rabbitmq/amqp091-go"

	"github.com/danblok/cryptoquotes/internal/logging"
	"github.com/danblok/cryptoquotes/pkg/config"
)

// MessageQueue is an interface for messaging queue.
type MessageQueue interface {
	EmittingChannel() (*mq.Channel, error)
	EmittingChannelTLS() (*mq.Channel, error)
	ReceivivgChannelAndQueue() (*mq.Channel, *mq.Queue, error)
	ReceivivgChannelAndQueueTLS() (*mq.Channel, *mq.Queue, error)
	LogMessages() error
	Close() error
}

type logMessageQueue struct {
	conn   *mq.Connection
	ch     *mq.Channel
	q      *mq.Queue
	config *config.AMQPConfig
}

// New returns a new logMessageQueue.
func New(config *config.AMQPConfig) MessageQueue {
	return &logMessageQueue{config: config}
}

func (l *logMessageQueue) EmittingChannel() (*mq.Channel, error) {
	conn, err := mq.Dial(l.config.URL())
	if err != nil {
		return nil, err
	}
	l.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	l.ch = ch

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return l.ch, nil
}

func (l *logMessageQueue) EmittingChannelTLS() (*mq.Channel, error) {
	conn, err := mq.DialTLS(
		l.config.URL(),
		l.config.TLSConfig(),
	)
	if err != nil {
		return nil, err
	}
	l.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	l.ch = ch

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return l.ch, nil
}

func (l *logMessageQueue) ReceivivgChannelAndQueue() (*mq.Channel, *mq.Queue, error) {
	conn, err := mq.Dial(l.config.URL())
	if err != nil {
		return nil, nil, err
	}
	l.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	l.ch = ch

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, nil, err
	}
	l.q = &q

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	if err != nil {
		return nil, nil, err
	}

	return l.ch, l.q, nil
}

func (l *logMessageQueue) ReceivivgChannelAndQueueTLS() (*mq.Channel, *mq.Queue, error) {
	conn, err := mq.DialTLS(l.config.URL(), l.config.TLSConfig())
	if err != nil {
		return nil, nil, err
	}
	l.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	l.ch = ch

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, nil, err
	}
	l.q = &q

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	if err != nil {
		return nil, nil, err
	}

	return l.ch, l.q, nil
}

func (l *logMessageQueue) LogMessages() error {
	msgs, err := l.ch.Consume(l.q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for m := range msgs {
		var log logging.Log
		err := gob.NewDecoder(bytes.NewReader(m.Body)).Decode(&log)
		if err != nil {
			slog.Error("decoding log message", "error", err)
			continue
		}

		if log.Error != "" {
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"request trace",
				slog.String("req_id", log.RequestID),
				slog.Time("start_t", log.StartTime),
				slog.Duration("req_time", log.RequestTime),
				slog.String("transport", log.TransportType),
				slog.String("remote_addr", log.RemoteAddr),
				slog.String("error", log.Error),
			)
		} else {
			slog.LogAttrs(
				context.Background(),
				slog.LevelInfo,
				"request trace",
				slog.String("req_id", log.RequestID),
				slog.Time("start_t", log.StartTime),
				slog.Duration("req_time", log.RequestTime),
				slog.String("transport", log.TransportType),
				slog.String("remote_addr", log.RemoteAddr),
				slog.Any("quotes", log.Quotes),
			)
		}
	}
	return nil
}

func (l *logMessageQueue) Close() error {
	if err := l.ch.Close(); err != nil {
		return err
	}

	return l.conn.Close()
}

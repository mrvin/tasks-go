package nats

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"

	"github.com/mrvin/tasks-go/goods/internal/storage"
	"github.com/nats-io/nats.go"
)

const EventsBuffer = 1000

type Conf struct {
	Host    string
	Port    string
	Subject string
}

type Queue struct {
	conn     *nats.Conn
	sub      *nats.Subscription
	subject  string
	closeCh  chan struct{}
	EventsCh chan storage.Event
}

func New(conf *Conf) (*Queue, error) {
	url := "nats://" + net.JoinHostPort(conf.Host, conf.Port)
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to queue: %w", err)
	}

	eventsCh := make(chan storage.Event, EventsBuffer)
	closeCh := make(chan struct{})
	sub, err := conn.Subscribe(conf.Subject, func(msg *nats.Msg) {
		var event storage.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			slog.Error("Unmarshal event: " + err.Error())
			return
		}
		select {
		case <-closeCh:
			return
		default:
			eventsCh <- event
		}
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}
	nats.DeliverNew()

	return &Queue{
		conn:     conn,
		sub:      sub,
		subject:  conf.Subject,
		closeCh:  closeCh,
		EventsCh: eventsCh,
	}, nil
}

func (q *Queue) SendEvent(event *storage.Event) error {
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if err := q.conn.Publish(q.subject, jsonEvent); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}

func (q *Queue) Close() error {
	close(q.closeCh)

	if err := q.sub.Unsubscribe(); err != nil {
		return fmt.Errorf("unsubscribing: %w", err)
	}
	close(q.EventsCh)
	q.conn.Close()

	return nil
}

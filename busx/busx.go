package busx

import (
	"context"
	"time"
)

type Conn interface {
	Close() error
}

type TopologyInspector interface {
	CheckQueues(names ...string) error
}

type Publisher interface {
	PublishToQueue(ctx context.Context, queueName string, payload []byte, opts ...PubOpt) error
	Close() error
}

type Consumer interface {
	ConsumeQueue(ctx context.Context, queueName string, handler Handler, opts ...ConOpt) error
	Close() error
}

type Handler func(ctx context.Context, d Delivery) Ack

type Ack int

const (
	AckOk Ack = iota
	AckNack
	AckRequeue
)

type Delivery struct {
	Body       []byte
	RoutingKey string
	Headers    map[string]any
}

type PubOpt func(*PubOptions)
type ConOpt func(*ConOptions)

type PubOptions struct {
	Mandatory bool
	Headers   map[string]any
}
type ConOptions struct {
	AutoAck bool
}

// Observabilidad mínima
type Observer interface {
	OnConnected(url string)
	OnDisconnected(err error)
	OnPublishToQueue(queue string, size int, err error)
	OnPublishToExchange(exchange, routingKey string, size int, err error)
	OnConsume(queue string, ack Ack, dur time.Duration, err error)
}

type NopObserver struct{}

func (NopObserver) OnConnected(string)                             {}
func (NopObserver) OnDisconnected(error)                           {}
func (NopObserver) OnPublishToQueue(string, int, error)            {}
func (NopObserver) OnPublishToExchange(string, string, int, error) {}
func (NopObserver) OnConsume(string, Ack, time.Duration, error)    {}

type Bundle struct {
	Conn      Conn
	Publisher Publisher
	Consumer  Consumer
	Inspector TopologyInspector
}

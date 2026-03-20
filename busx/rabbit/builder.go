package rabbit

import (
	"crypto/tls"
	"fmt"
	"github.com/Mozart-SymphonIA/infra-mz/busx"

	amqp "github.com/rabbitmq/amqp091-go"
)

func BuildRabbit(cfg busx.Config) (*busx.Bundle, error) {
	if cfg.Observer == nil {
		cfg.Observer = busx.NopObserver{}
	}
	dialCfg := amqp.Config{
		Heartbeat:       cfg.Heartbeat,
		Locale:          cfg.Locale,
		TLSClientConfig: tlsIfAmqps(cfg.URL),
	}
	conn, err := amqp.DialConfig(cfg.URL, dialCfg)
	if err != nil {
		cfg.Observer.OnDisconnected(err)
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}
	cfg.Observer.OnConnected(cfg.URL)

	pubCh, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("publisher channel: %w", err)
	}
	var confCh <-chan amqp.Confirmation
	if cfg.PublisherConfirms {
		if err := pubCh.Confirm(false); err != nil {
			_ = pubCh.Close()
			_ = conn.Close()
			return nil, fmt.Errorf("publisher confirm: %w", err)
		}
		confCh = pubCh.NotifyPublish(make(chan amqp.Confirmation, 1024))
	}

	conCh, err := conn.Channel()
	if err != nil {
		_ = pubCh.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("consumer channel: %w", err)
	}
	prefetch := cfg.Prefetch
	if prefetch <= 0 {
		prefetch = 1
	}
	if err := conCh.Qos(prefetch, 0, false); err != nil {
		_ = conCh.Close()
		_ = pubCh.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("consumer qos: %w", err)
	}

	return &busx.Bundle{
		Conn:      &rabbitConn{c: conn, ob: cfg.Observer},
		Publisher: &rabbitPublisher{ch: pubCh, confirms: confCh, useConf: cfg.PublisherConfirms, ob: cfg.Observer},
		Consumer:  &rabbitConsumer{ch: conCh, ob: cfg.Observer},
		Inspector: &rabbitInspector{conn: conn},
	}, nil
}

func tlsIfAmqps(url string) *tls.Config {
	if len(url) >= 8 && url[:8] == "amqps://" {
		return &tls.Config{MinVersion: tls.VersionTLS12}
	}
	return nil
}

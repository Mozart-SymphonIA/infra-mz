package rabbit

import (
	"context"
	"fmt"
	"github.com/Mozart-SymphonIA/infra-mz/busx"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	_ busx.Consumer = (*rabbitConsumer)(nil)
)

type rabbitConsumer struct {
	ch *amqp.Channel
	ob busx.Observer
}

func (c *rabbitConsumer) Close() error { return c.ch.Close() }

func (c *rabbitConsumer) ConsumeQueue(ctx context.Context, queueName string, handler busx.Handler, opts ...busx.ConOpt) error {
	o := busx.ConOptions{AutoAck: false}
	for _, fn := range opts {
		fn(&o)
	}

	_, err := c.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue %s: %w", queueName, err)
	}

	msgs, err := c.ch.Consume(queueName, "", o.AutoAck, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case m, ok := <-msgs:
			if !ok {
				return nil
			}
			c.processMessage(ctx, queueName, m, handler)
		}
	}
}

func (c *rabbitConsumer) processMessage(ctx context.Context, queue string, m amqp.Delivery, handler busx.Handler) {
	start := time.Now()
	d := busx.Delivery{
		Body:       m.Body,
		RoutingKey: m.RoutingKey,
		Headers:    tableToMap(m.Headers),
	}
	ack := handler(ctx, d)
	switch ack {
	case busx.AckOk:
		_ = m.Ack(false)
	case busx.AckNack:
		_ = m.Nack(false, false)
	case busx.AckRequeue:
		_ = m.Nack(false, true)
	}
	c.ob.OnConsume(queue, ack, time.Since(start), nil)
}

func tableToMap(t amqp.Table) map[string]any {
	if t == nil {
		return nil
	}
	out := make(map[string]any, len(t))
	for k, v := range t {
		out[k] = v
	}
	return out
}

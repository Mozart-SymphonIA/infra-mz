package rabbit

import (
	"context"
	"fmt"
	"github.com/Mozart-SymphonIA/infra-mz/busx"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	_ busx.Publisher = (*rabbitPublisher)(nil)
)

type rabbitPublisher struct {
	ch       *amqp.Channel
	confirms <-chan amqp.Confirmation
	useConf  bool
	ob       busx.Observer
}

func (p *rabbitPublisher) Close() error { return p.ch.Close() }

func (p *rabbitPublisher) PublishToQueue(ctx context.Context, queueName string, payload []byte, opts ...busx.PubOpt) error {
	o := busx.PubOptions{}
	for _, fn := range opts {
		fn(&o)
	}

	pub := amqp.Publishing{
		Headers:      amqp.Table{},
		ContentType:  "application/octet-stream",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         payload,
	}
	for k, v := range o.Headers {
		pub.Headers[k] = v
	}

	if err := p.ch.PublishWithContext(ctx, "", queueName, o.Mandatory, false, pub); err != nil {
		p.ob.OnPublishToQueue(queueName, len(payload), err)
		return err
	}
	if p.useConf {
		select {
		case conf, ok := <-p.confirms:
			if !ok || !conf.Ack {
				err := fmt.Errorf("publish to queue not confirmed (ack=%v ok=%v)", conf.Ack, ok)
				p.ob.OnPublishToQueue(queueName, len(payload), err)
				return err
			}
		case <-ctx.Done():
			err := ctx.Err()
			p.ob.OnPublishToQueue(queueName, len(payload), err)
			return err
		}
	}
	p.ob.OnPublishToQueue(queueName, len(payload), nil)
	return nil
}

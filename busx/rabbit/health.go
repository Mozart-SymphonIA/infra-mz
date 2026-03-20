package rabbit

import (
	"github.com/Mozart-SymphonIA/infra-mz/busx"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitInspector struct {
	conn *amqp.Connection
}

var (
	_ busx.TopologyInspector = (*rabbitInspector)(nil)
)

func (ri *rabbitInspector) CheckQueues(names ...string) error {
	ch, err := ri.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	for _, q := range names {
		if _, err := ch.QueueDeclarePassive(q, true, false, false, false, nil); err != nil {
			return err
		}
	}
	return nil
}

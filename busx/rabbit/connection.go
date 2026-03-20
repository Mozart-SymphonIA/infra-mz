package rabbit

import (
	"github.com/Mozart-SymphonIA/infra-mz/busx"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	_ busx.Conn = (*rabbitConn)(nil)
)

type rabbitConn struct {
	c  *amqp.Connection
	ob busx.Observer
}

func (r *rabbitConn) Close() error { return r.c.Close() }

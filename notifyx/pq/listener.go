package pq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	libpq "github.com/lib/pq"

	"github.com/Mozart-SymphonIA/infra-mz/notifyx"
)

// listener implementa notifyx.Listener usando lib/pq.
// lib/pq gestiona la reconexión automática con backoff configurable.
type listener struct {
	l *libpq.Listener
}

func New(connStr string) notifyx.Listener {
	l := libpq.NewListener(connStr, 5*time.Second, time.Minute, func(ev libpq.ListenerEventType, err error) {
		if err != nil {
			slog.Error("notifyx: reconnect", "event", ev, "err", err)
		}
	})
	return &listener{l: l}
}

func (l *listener) Listen(ctx context.Context, channel string, handler notifyx.Handler) error {
	if err := l.l.Listen(channel); err != nil {
		return fmt.Errorf("notifyx: listen %q: %w", channel, err)
	}
	defer l.l.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		case notif, ok := <-l.l.NotificationChannel():
			if !ok {
				return fmt.Errorf("notifyx: notification channel closed")
			}
			if notif == nil {
				// ping periódico de lib/pq para mantener la conexión viva
				continue
			}
			handler(ctx, notif.Extra)
		}
	}
}

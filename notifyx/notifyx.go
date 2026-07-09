package notifyx

import "context"

// Handler es la función que se invoca cuando llega una notificación.
// payload es el string libre que envió el emisor via pg_notify.
type Handler func(ctx context.Context, payload string)

// Listener se suscribe a un canal de Postgres y entrega notificaciones al Handler.
// La implementación maneja la reconexión automática ante caídas de red.
type Listener interface {
	Listen(ctx context.Context, channel string, handler Handler) error
}

// Publisher emite notificaciones a un canal de Postgres via pg_notify.
type Publisher interface {
	Publish(ctx context.Context, channel, payload string) error
}

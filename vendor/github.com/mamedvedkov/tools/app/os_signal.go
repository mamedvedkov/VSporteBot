package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ErrorOSSignal обертка в виде ошибки над сигналом ОС для того, чтобы использовать SignalNotify в errgroup.
type ErrorOSSignal struct {
	signal os.Signal
}

func (e *ErrorOSSignal) Error() string {
	return e.signal.String()
}

// SignalNotify обработка событий выхода от ОС.
func SignalNotify(ctx context.Context) error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-c:
		return &ErrorOSSignal{signal: s}
	case <-ctx.Done():
		return ctx.Err()
	}
}

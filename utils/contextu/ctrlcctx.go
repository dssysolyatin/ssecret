package contextu

import (
	"context"
	"os"
	"os/signal"
)

func WithCancelCtrlC(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case <-c:
		}

		cancel()
		signal.Stop(c)
	}()

	return ctx
}

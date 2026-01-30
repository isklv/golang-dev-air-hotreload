package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/isklv/slogging"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	slogging.L(ctx).Info("Hello world")

	<-ctx.Done()

	slogging.L(ctx).Info("Shutting down")
}

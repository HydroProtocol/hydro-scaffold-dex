package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WaitExitSignal(ctxStop context.CancelFunc) {
	var exitSignal = make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGTERM)
	signal.Notify(exitSignal, syscall.SIGINT)

	sig := <-exitSignal
	println("caught sig: %+v, Stopping...", sig)
	ctxStop()
}

package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(cancel context.CancelFunc) {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	cancel()
}

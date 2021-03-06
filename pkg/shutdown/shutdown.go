package shutdown

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown makes sure that all resources are disposed correctly
func GracefulShutdown(cancel context.CancelFunc, logger *logrus.Logger) {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	logger.Info("Shutting down application.")
	cancel()
}

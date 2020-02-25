package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func NewSignalMonitoringChannel() chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	return signalChannel
}

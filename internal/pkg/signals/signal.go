package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// Create os.Signal channel that will trigger if specified interrupts are applied to process.
func NewSignalMonitoringChannel() chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	return signalChannel
}

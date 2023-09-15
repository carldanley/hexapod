package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/carldanley/hexapod/pkg/servos"
)

var signalChannel chan os.Signal

func init() {
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)
}

func main() {
	var testServoType = servos.NewServoType(135, 538, 0, 270.0, -90.0, 90.0)
	testServoType.DebugOutput()

	// <-signalChannel
}

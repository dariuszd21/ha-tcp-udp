package main

import (
	// "ha-tcp-udp/config"
	"ha-tcp-udp/logger"
	// "ha-tcp-udp/tcp_server"
)

func main() {

	logger.Debug("Some debug print")
	logger.SetLogLevel(logger.DEBUG)
	logger.Debug("Some debug print")
}

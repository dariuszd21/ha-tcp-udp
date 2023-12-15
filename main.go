package main

import (
	// "ha-tcp-udp/config"
	"ha-tcp-udp/logger"
	// "ha-tcp-udp/tcp_server"
)

func main() {
	logger_obj := logger.NewLogger(true)

	logger_obj.Debug("Some debug print")
}

package main

import (
	// "ha-tcp-udp/config"
	"ha-tcp-udp/logger"
	"ha-tcp-udp/tcp_server"
)

func main() {

	logger.SetLogLevel(logger.DEBUG)
	logger.Debug("Some debug print")

	server_config := tcp_server.ServerConfig{
		TcpHost: "localhost",
		TcpPort: 12000,
	}
	server := tcp_server.CreateServer(&server_config)
	server.Bind()
	server.Serve()
}

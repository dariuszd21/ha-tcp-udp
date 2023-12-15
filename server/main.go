package main

import (
	// "ha-tcp-udp/config"
	"ha-tcp-udp/logger"
	"ha-tcp-udp/server_if"
	"ha-tcp-udp/tcp_server"
	"ha-tcp-udp/udp_server"
	"os"
	"os/signal"
	"syscall"
)

func runServer(server server_if.Server) {
	server.Bind()
	server.Serve()
}

func main() {
	logger.SetLogPrint(true)
	logger.SetLogLevel(logger.DEBUG)

	server_config := server_if.ServerConfig{
		Host: "localhost",
		Port: 12000,
	}
	server := tcp_server.CreateServer(&server_config)
	go runServer(server)
	udp_server_config := server_if.ServerConfig{
		Host: "localhost",
		Port: 13000,
	}
	udp_server := udp_server.CreateServer(&udp_server_config)
	go runServer(udp_server)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		logger.Debugf("Received signal: %d", sig)
		done <- true
	}()

	<-done
	logger.Debug("exiting")
}

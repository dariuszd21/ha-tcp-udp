package main

import (
	"flag"
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
	tcp_port := flag.Int("tcp_port", 12000, "Port on which TCP server listen")
	tcp_host := flag.String("tcp_host", "localhost", "Address on which TCP server listen")
	udp_port := flag.Int("udp_port", 13000, "Port on which UDP server listen")
	udp_host := flag.String("udp_host", "localhost", "Address on which UDP server listen")
	print_logs := flag.Bool("print_logs", false, "Check if logs should be printed using PrintLn")
	log_level := flag.Int("log_level", logger.NONE, "Log level: 0 - None, 1 - Error, 2 - Debug")

	flag.Parse()

	logger.SetLogPrint(*print_logs)
	logger.SetLogLevel(*log_level)

	server_config := server_if.ServerConfig{
		Host: *tcp_host,
		Port: *tcp_port,
	}
	server := tcp_server.CreateServer(&server_config)
	go runServer(server)
	udp_server_config := server_if.ServerConfig{
		Host: *udp_host,
		Port: *udp_port,
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

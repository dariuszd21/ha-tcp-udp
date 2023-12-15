package tcp_server

import (
	"encoding/json"
	"fmt"
	"ha-tcp-udp/logger"
	"net"
	"sync/atomic"
	"time"
)

type ServerConfig struct {
	tcp_host          string
	tcp_port          int
	connections_limit uint64 // 0 means no-limit
}

type TCPServer struct {
	config      *ServerConfig
	listener    *net.Listener
	connections atomic.Uint64
}

func (server *TCPServer) bind() {
	listener, error := net.Listen("tcp",
		fmt.Sprintf("%s:%d", server.config.tcp_host, server.config.tcp_port),
	)
	if error != nil {
		logger.Fatal(error.Error())
	}
	server.listener = &listener
}

func (server *TCPServer) serve() {
	if server.listener == nil {
		logger.Fatal("Cannot serve connections. bind() forgetten?")
	}

	defer (*server.listener).Close()

	for {
		// TODO: Add handling sigterm

		if server.config.connections_limit != 0 {
			// Check if connection limit is already reached
			if server.connections.Load() >= server.config.connections_limit {
				logger.Debug("Server no longer accepts connections.")
				time.Sleep(time.Second)
			}
		}
		conn, err := (*server.listener).Accept()
		if err != nil {
			logger.Errorf("Cannot accept connection: %s", err.Error())
		}
		go server.handleConnection(conn)
	}
}

func (server *TCPServer) handleConnection(connection net.Conn) {
	logger.Debug("Connection opened")
	server.connections.Add(1)
	defer server.closeConnection(connection)
	buf := make([]byte, 1024)
	// TODO: Add timeout
	read_bytes, err := connection.Read(buf)

	if err != nil {
		logger.Errorf("Cannot handle connection: %s", connection.RemoteAddr())
		return
	}

	var json_response map[string]any

	if read_bytes == 0 {
		logger.Error("Bytes not read. Expected message")
		return
	}

	err = json.Unmarshal(buf, &json_response)

	if err != nil {
		logger.Error("Cannot decode bytes")
		logger.Debugf("%s", buf)
		return
	}

	// TODO: Add
	i := 0
	for {
		write_buff, err := json.Marshal(map[string]int{"response": i})
		if err != nil {
			logger.Errorf("Cannot encode response: %d", i)
			return
		}
		_, err = connection.Write(write_buff)
		if err != nil {
			logger.Errorf("Cannot write response: %s", err)
			return
		}

		time.Sleep(50 * time.Millisecond)
		i++
	}
}

func (server *TCPServer) closeConnection(connection net.Conn) {
	logger.Debug("Closing connection")
	server.connections.Add(^(uint64(0)))
	connection.Close()
}

func CreateServer(server_config *ServerConfig) *TCPServer {
	return &TCPServer{
		config: server_config,
	}
}

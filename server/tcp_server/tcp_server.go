package tcp_server

import (
	"encoding/json"
	"fmt"
	"ha-tcp-udp/logger"
	"ha-tcp-udp/server_if"
	"net"
	"sync/atomic"
	"time"
)

type TCPServer struct {
	config      *server_if.ServerConfig
	listener    net.Listener
	connections atomic.Uint64
}

func (server *TCPServer) Bind() {
	listener, error := net.Listen("tcp",
		fmt.Sprintf("%s:%d", server.config.Host, server.config.Port),
	)
	if error != nil {
		logger.Fatal(error.Error())
	}
	logger.Debugf("Starting TCP server on %s:%d", server.config.Host, server.config.Port)
	server.listener = listener
}

func (server *TCPServer) Serve() {
	if server.listener == nil {
		logger.Fatal("Cannot serve connections. bind() forgetten?")
	}

	defer server.listener.Close()

	session_id_chan := make(chan server_if.SessionIdOps)
	go server_if.AssignSessionId(session_id_chan)

	for {
		if server.config.ConnectionsLimit != 0 {
			// Check if connection limit is already reached
			if server.connections.Load() >= server.config.ConnectionsLimit {
				logger.Debug("Server no longer accepts connections.")
				time.Sleep(time.Second)
			}
		}
		conn, err := server.listener.Accept()
		if err != nil {
			logger.Errorf("Cannot accept connection: %s", err.Error())
		}
		go server.handleConnection(conn, session_id_chan)
	}
}

func (server *TCPServer) handleConnection(connection net.Conn, assignIdChannel chan server_if.SessionIdOps) {
	const DELAY_BETWEEN_MESAGES = 50 * time.Millisecond
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

	if read_bytes == 0 {
		logger.Error("Bytes not read. Expected message")
		return
	}

	var init_message server_if.InitMessage
	err = json.Unmarshal(buf[:read_bytes], &init_message)

	if err != nil {
		logger.Error("Cannot decode bytes")
		logger.Error(err.Error())
		logger.Debugf("%s", buf)
		return
	}

	if init_message.SessionId != 0 {
		logger.Debugf("Session: %d reestablished", init_message.SessionId)
	} else {
		session_id_resp := make(chan uint64)
		session_id_ops := server_if.SessionIdOps{
			Resp: session_id_resp,
		}
		assignIdChannel <- session_id_ops
		session_id := <-session_id_resp

		session_est, err := json.Marshal(&server_if.SessionEstablishment{SessionId: session_id})
		if err != nil {
			logger.Errorf("Cannot encode session establishment: %s", session_est)
			return
		}
		_, err = connection.Write(session_est)
		if err != nil {
			logger.Error("Cannot send session establishment message")
			logger.Debug(err.Error())
			return
		}
		time.Sleep(DELAY_BETWEEN_MESAGES)
	}

	i := init_message.LastMessage + 1
	for {
		write_buff, err := json.Marshal(map[string]uint64{"response": i})
		if err != nil {
			logger.Errorf("Cannot encode response: %d", i)
			return
		}
		_, err = connection.Write(write_buff)
		if err != nil {
			logger.Errorf("Cannot send response: %s", err)
			return
		}

		time.Sleep(DELAY_BETWEEN_MESAGES)
		i++
	}
}

func (server *TCPServer) closeConnection(connection net.Conn) {
	logger.Debug("Closing connection")
	server.connections.Add(^(uint64(0)))
	connection.Close()
}

func CreateServer(server_config *server_if.ServerConfig) *TCPServer {
	return &TCPServer{
		config: server_config,
	}
}

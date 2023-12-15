package udp_server

import (
	"encoding/json"
	"fmt"
	"ha-tcp-udp/logger"
	"ha-tcp-udp/server_if"
	"net"
	"sync/atomic"
	"time"
)

type UDPServer struct {
	config      *server_if.ServerConfig
	listener    *net.PacketConn
	connections atomic.Uint64
}

func (server *UDPServer) Bind() {
	listener, error := net.ListenPacket("udp",
		fmt.Sprintf("%s:%d", server.config.Host, server.config.Port),
	)
	if error != nil {
		logger.Fatal(error.Error())
	}
	logger.Debugf("Starting UDP server on %s:%d", server.config.Host, server.config.Port)
	server.listener = &listener

}

func (server *UDPServer) Serve() {
	if server.listener == nil {
		logger.Fatal("Cannot serve connections. bind() forgetten?")
	}

	defer (*server.listener).Close()

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
		buf := make([]byte, 1024)
		n_bytes, addr, err := (*server.listener).ReadFrom(buf)
		if err != nil {
			logger.Errorf("Cannot read from: %s", err.Error())
			return
		}
		init_message, err := readInitMessage(buf, n_bytes)
		go server.handleConnection(addr, init_message, session_id_chan)
	}
}

func readInitMessage(buff []byte, n_bytes int) (*server_if.InitMessage, error) {

	if n_bytes == 0 {
		logger.Error("Bytes not read. Expected message")
		return nil, nil
	}

	var init_message server_if.InitMessage
	err := json.Unmarshal(buff[:n_bytes], &init_message)

	if err != nil {
		logger.Error("Cannot decode bytes")
		logger.Error(err.Error())
		logger.Debugf("%s", buff)
		return nil, err
	}

	return &init_message, nil
}

func (server *UDPServer) handleConnection(addr net.Addr, init_message *server_if.InitMessage, assignIdChannel chan server_if.SessionIdOps) {
	logger.Debug("Connection opened")
	server.connections.Add(1)
	defer server.closeConnection()

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
		_, err = (*server.listener).WriteTo(session_est, addr)
		if err != nil {
			logger.Error("Cannot send session establishment message")
			logger.Debug(err.Error())
			return
		}
	}

	i := init_message.LastMessage + 1
	for {
		write_buff, err := json.Marshal(map[string]uint64{"response": i})
		if err != nil {
			logger.Errorf("Cannot encode response: %d", i)
			return
		}
		_, err = (*server.listener).WriteTo(write_buff, addr)
		if err != nil {
			logger.Errorf("Cannot send response: %s", err)
			return
		}

		time.Sleep(50 * time.Millisecond)
		i++
	}
}

func (server *UDPServer) closeConnection() {
	logger.Debug("Closing connection")
	server.connections.Add(^(uint64(0)))
}

func CreateServer(server_config *server_if.ServerConfig) *UDPServer {
	return &UDPServer{
		config: server_config,
	}
}

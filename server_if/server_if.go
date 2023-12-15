package server_if

import (
	"ha-tcp-udp/logger"
)

type InitMessage struct {
	SessionId   uint64 `json:"session_id"`
	LastMessage uint64 `json:"last_message"`
}

type SessionEstablishment struct {
	SessionId uint64 `json:"session_id"`
}

type SessionIdOps struct {
	Resp chan uint64
}

func AssignSessionId(req chan SessionIdOps) {
	var i uint64 = 1

	for {
		request := <-req
		request.Resp <- i
		logger.Debugf("Assigning session id: %d", i)
		i++
		if i == 0 {
			i++
		}
	}
}

type Server interface {
	Bind()
	Serve()
}

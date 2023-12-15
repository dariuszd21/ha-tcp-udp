package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type InitMessage struct {
	SessionId   uint64 `json:"session_id"`
	LastMessage uint64 `json:"last_message"`
}

type SessionEstablishment struct {
	SessionId uint64 `json:"session_id"`
}

func runStable(host string, port int, connNumber int) {
	for i := 0; i < connNumber; i++ {
		i := i
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				log.Println("Error:", err)
				return
			}
			log.Println("Connection established", i)
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				log.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				log.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				log.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				log.Println("Cannot read init message", err)
				log.Println(string(buff))
				return
			}
			for {
				_, err = conn.Read(buff)
				if err != nil {
					log.Println("Cannot read from server", err)
					return
				}
			}
		}()
	}
}

func runReconnecting(host string, port int, connNumber int) {
	for i := 0; i < connNumber; i++ {
		i := i
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			random_retry := 50 + rand.Intn(1000)
			log.Printf("Client will retry every %d packets\n", random_retry)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			log.Println("Connection established", i)
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				log.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				log.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				log.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				log.Println("Cannot read init message", err)
				return
			}
			for read_nr := 1; ; read_nr++ {
				if read_nr%random_retry == 0 {
					log.Println("Reestablished on", read_nr)
					conn.Close()
					conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
					if err != nil {
						log.Println("Error:", err)
						return
					}
					defer conn.Close()
					init_msg := InitMessage{
						SessionId:   sessionEstablishment.SessionId,
						LastMessage: uint64(read_nr),
					}
					handshake_buff, err := json.Marshal(&init_msg)
					if err != nil {
						log.Println("Cannot encode init message")
						return
					}
					_, err = conn.Write(handshake_buff)
					if err != nil {
						log.Println("Cannot send init message")
						return
					}
				}

				_, err = conn.Read(buff)
				if err != nil {
					log.Println("Cannot read from server", err)
					return
				}
			}
		}()
	}
}

func runDropping(host string, port int, connNumber int) {
	done_channel := make(chan bool)
	for i := 0; i < connNumber; i++ {
		i := i
		go func() {
			random_retry := 50 + rand.Intn(1000)
			log.Printf("Client will drop after %d packets\n", random_retry)
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				log.Println("Error:", err)
				return
			}
			log.Println("Connection established", i)
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				log.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				log.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				log.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				log.Println("Cannot read init message", err)
				log.Println(string(buff))
				return
			}
			for packet_nr := 0; packet_nr < random_retry; packet_nr++ {
				_, err = conn.Read(buff)
				if err != nil {
					log.Println("Cannot read from server", err)
					return
				}
			}
			done_channel <- true
		}()
	}

	for {
		<-done_channel
		go func() {
			random_retry := 50 + rand.Intn(1000)
			log.Printf("Client will drop after %d packets\n", random_retry)
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				log.Println("Error:", err)
				return
			}
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				log.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				log.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				log.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				log.Println("Cannot read init message", err)
				log.Println(string(buff))
				return
			}
			for packet_nr := 0; packet_nr < random_retry; packet_nr++ {
				_, err = conn.Read(buff)
				if err != nil {
					log.Println("Cannot read from server", err)
					return
				}
			}
			done_channel <- true
		}()

	}
}

func main() {
	stable := flag.Int("stable", 0, "number of stable connections")
	reconnecting := flag.Int("reconnecting", 0, "number of reconnecting connections")
	dropping := flag.Int("dropping", 0, "number of dropped connections")
	host := flag.String("host", "localhost", "host on which server is listening")
	port := flag.Int("port", 12000, "port on which server is listening")

	flag.Parse()

	log.Println("Stable connections", *stable)
	log.Println("Reconnecting connections", *reconnecting)
	log.Println("Dropping connections", *dropping)
	log.Printf("Connection host %s:%d\n", *host, *port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go runStable(*host, *port, *stable)
	go runReconnecting(*host, *port, *reconnecting)
	go runDropping(*host, *port, *dropping)

	go func() {
		sig := <-sigs
		fmt.Printf("Received signal: %d\n", sig)
		done <- true
	}()

	<-done
	log.Println("exiting...")
}

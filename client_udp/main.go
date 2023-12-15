package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func udp_connect(host string, port int) (*net.UDPConn, error) {
	remote_add, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	local_add, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return net.DialUDP("udp", local_add, remote_add)
}

func runStable(host string, port int, connNumber int) {
	for i := 0; i < connNumber; i++ {
		i := i
		go func() {
			conn, err := udp_connect(host, port)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Connection established", i)
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				fmt.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				fmt.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				fmt.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				fmt.Println("Cannot read init message", err)
				fmt.Println(string(buff))
				return
			}
			for {
				_, err = conn.Read(buff)
				if err != nil {
					fmt.Println("Cannot read from server", err)
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
			conn, err := udp_connect(host, port)
			random_retry := 50 + rand.Intn(1000)
			fmt.Printf("Client will retry every %d packets\n", random_retry)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Connection established", i)
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				fmt.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				fmt.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				fmt.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				fmt.Println("Cannot read init message", err)
				return
			}
			for read_nr := 1; ; read_nr++ {
				if read_nr%random_retry == 0 {
					fmt.Println("Reestablished on", read_nr)
					conn, err := udp_connect(host, port)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					defer conn.Close()
					init_msg := InitMessage{
						SessionId:   sessionEstablishment.SessionId,
						LastMessage: uint64(read_nr),
					}
					handshake_buff, err := json.Marshal(&init_msg)
					if err != nil {
						fmt.Println("Cannot encode init message")
						return
					}
					_, err = conn.Write(handshake_buff)
					if err != nil {
						fmt.Println("Cannot send init message")
						return
					}
				}

				_, err = conn.Read(buff)
				if err != nil {
					fmt.Println("Cannot read from server", err)
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
			fmt.Printf("Client will drop after %d packets\n", random_retry)
			conn, err := udp_connect(host, port)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Connection established", i)
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				fmt.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				fmt.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				fmt.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				fmt.Println("Cannot read init message", err)
				fmt.Println(string(buff))
				return
			}
			for packet_nr := 0; packet_nr < random_retry; packet_nr++ {
				_, err = conn.Read(buff)
				if err != nil {
					fmt.Println("Cannot read from server", err)
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
			fmt.Printf("Client will drop after %d packets\n", random_retry)
			conn, err := udp_connect(host, port)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				fmt.Println("Cannot encode init message")
				return
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				fmt.Println("Cannot send init message")
				return
			}

			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				fmt.Println("Cannot read session establishment")
				return
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				fmt.Println("Cannot read init message", err)
				fmt.Println(string(buff))
				return
			}
			for packet_nr := 0; packet_nr < random_retry; packet_nr++ {
				_, err = conn.Read(buff)
				if err != nil {
					fmt.Println("Cannot read from server", err)
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
	port := flag.Int("port", 13000, "port on which server is listening")

	flag.Parse()

	fmt.Println("Stable connections", *stable)
	fmt.Println("Reconnecting connections", *reconnecting)
	fmt.Println("Dropping connections", *dropping)
	fmt.Printf("Connection host %s:%d\n", *host, *port)

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
	fmt.Println("exiting...")
}

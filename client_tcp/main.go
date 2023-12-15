package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		go func ()  {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    		if err != nil {
        		fmt.Println("Error:", err)
        		return
    		}
			fmt.Println("Connection established", i)
    		defer conn.Close()
			handshake_buff, err := json.Marshal(&InitMessage{})
			if err != nil {
				fmt.Println("Cannot encode init message")
			}
			_, err = conn.Write(handshake_buff)
			if err != nil {
				fmt.Println("Cannot send init message")
			}
			
			buff := make([]byte, 1024)
			n_bytes, err := conn.Read(buff)
			if err != nil {
				fmt.Println("Cannot read session establishment")
			}
			var sessionEstablishment SessionEstablishment
			err = json.Unmarshal(buff[:n_bytes], &sessionEstablishment)
			if err != nil {
				fmt.Println("Cannot read init message", err)
				return
			}
			for {
				_, err = conn.Read(buff)
				if err != nil {
					fmt.Println("Cannot read from server", err)
				}
			}
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

	fmt.Println("Stable connections", *stable)
	fmt.Println("Reconnecting connections", *reconnecting)
	fmt.Println("Dropping connections", *dropping)
	fmt.Printf("Connection host %s:%d\n", *host, *port)

	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    done := make(chan bool, 1)

	runStable(*host, *port, *stable)

    go func() {
        sig := <-sigs
    	fmt.Printf("Received signal: %d\n", sig)
        done <- true
    }()

    <-done
    fmt.Println("exiting...")
}

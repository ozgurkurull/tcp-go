package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var currentConnectionCount int = 0

//https://go.dev/blog/pipelines
//https://stackoverflow.com/questions/25306073/always-have-x-number-of-goroutines-running-at-any-time

func main() {
	ln, err := net.Listen("tcp", ":10001")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	log.Printf("listen started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("accept error: %v", err)
		}
		currentConnectionCount++
		go serve(conn)
		log.Printf("Accepted %d", currentConnectionCount)
	}
}

func serve(conn net.Conn) {
	bufr := bufio.NewReader(conn)
	for {
		line, err := bufr.ReadString('\n')
		if err != nil {
			currentConnectionCount--
			return
		}

		log.Printf("Data: %s", line)

		response := fmt.Sprintf("OK : %s", line)
		conn.Write([]byte(response))
	}
}

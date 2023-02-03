package server2

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
)

var total_rcv int64

func main() {

	var cmd_rate_int float64
	var cmd_port string
	var client_size int

	flag.Float64Var(&cmd_rate_int, "rate", 400000, "change rate of message reading")
	flag.StringVar(&cmd_port, "port", ":9090", "port to listen")
	flag.IntVar(&client_size, "size", 20, "number of clients")

	flag.Parse()

	t := flag.Arg(0)

	if t == "server" {
		server(cmd_port)

	} else if t == "client" {
		for i := 0; i < client_size; i++ {
			go client(cmd_rate_int, cmd_port)
		}
		// <-make(chan bool) // infinite wait.
		<-time.After(time.Second * 2)
		fmt.Println("total exchanged", total_rcv)

	} else if t == "client_ratelimit" {
		bucket := ratelimit.NewBucketWithQuantum(time.Second, int64(cmd_rate_int), int64(cmd_rate_int))
		for i := 0; i < client_size; i++ {
			go clientRateLimite(bucket, cmd_port)
		}
		// <-make(chan bool) // infinite wait.
		<-time.After(time.Second * 3)
		fmt.Println("total exchanged", total_rcv)
	}
}

func server(cmd_port string) {
	ln, err := net.Listen("tcp", cmd_port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go io.Copy(conn, conn)
	}
}

func client(cmd_rate_int float64, cmd_port string) {

	conn, err := net.Dial("tcp", cmd_port)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func(conn net.Conn) {
		buf := make([]byte, 8)
		for {
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				break
			}
			// int_message := int64(binary.LittleEndian.Uint64(buf))
			// t2 := time.Unix(0, int_message)
			// fmt.Println("ROUDNTRIP", time.Now().Sub(t2))
			atomic.AddInt64(&total_rcv, 1)
		}
		return
	}(conn)

	byte_message := make([]byte, 8)
	for {
		wait := time.Microsecond * time.Duration(nextTime(cmd_rate_int))
		if wait > 0 {
			time.Sleep(wait)
			fmt.Println("WAIT", wait)
		}
		int_message := time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		_, err := conn.Write(byte_message)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
	}
}

func clientRateLimite(bucket *ratelimit.Bucket, cmd_port string) {

	conn, err := net.Dial("tcp", cmd_port)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func(conn net.Conn) {
		buf := make([]byte, 8)
		for {
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				break
			}
			// int_message := int64(binary.LittleEndian.Uint64(buf))
			// t2 := time.Unix(0, int_message)
			// fmt.Println("ROUDNTRIP", time.Now().Sub(t2))
			atomic.AddInt64(&total_rcv, 1)
		}
		return
	}(conn)

	byte_message := make([]byte, 8)
	for {
		bucket.Wait(1)
		int_message := time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		_, err := conn.Write(byte_message)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
	}
}

func nextTime(rate float64) float64 {
	return -1 * math.Log(1.0-rand.Float64()) / rate
}

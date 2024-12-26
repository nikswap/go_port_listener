package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	for port := 1024; port < 30000; port++ {
		go start_to_listen(":" + strconv.Itoa(port))
	}
	for {
	}
}

func start_to_listen(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listeing on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, port)
	}
}

func handleConnection(conn net.Conn, port string) {
	defer conn.Close()
	fmt.Println("DATA ON PORT", port)

	buf := make([]byte, 1024)
	timeoutDuration := 60 * time.Second
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	_, err := conn.Read(buf)
	lastNull := 0
	for idx := len(buf) - 1; idx >= 0; idx-- {
		if buf[idx] != 0x00 {
			break
		}
		lastNull = idx
	}
	buf = buf[:lastNull+1]
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Received: %s", base64.StdEncoding.EncodeToString(buf))
}

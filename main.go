package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		os.Exit(-1)
	}
}

func main() {
	startPort, err := strconv.Atoi(os.Args[1])
	CheckError(err)
	endPort, err := strconv.Atoi(os.Args[2])
	CheckError(err)
	for port := startPort; port < endPort; port++ {
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
	fmt.Printf("DATA ON PORT: %s FROM HOST: %s\n", port, conn.RemoteAddr().String())

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
	fmt.Printf("Received on %s: %s\n", port, base64.StdEncoding.EncodeToString(buf))
}

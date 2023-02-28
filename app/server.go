package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)

		}

		req := strings.Split(string(buff[:n]), "\r\n")
		command := req[2]
		if command == "echo" || command == "ECHO" {
			resp := fmt.Sprintf("+%s\r\n", req[4])
			conn.Write([]byte(resp))
		}
		if command == "ping" || command == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}

	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	storage := NewStorage()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleConn(conn, storage)
	}
}

func HandleConn(conn net.Conn, storage *InMemoryStorage) {
	defer conn.Close()
	for {
		value, err := DecodeRESP(bufio.NewReader(conn))

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}

		command := value.Array()[0].String()
		args := value.Array()[1:]

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
		case "set":
			if len(args) > 2 && args[2].String() == "px" {
				expiration, err := strconv.Atoi(args[3].String())
				if err != nil {
					conn.Write([]byte("-ERR\r\n"))
				} else {
					storage.Set(args[0].String(), args[1].String(), int64(expiration))
					conn.Write([]byte("+OK\r\n"))
				}
			} else {
				storage.Set(args[0].String(), args[1].String(), int64(0))
				conn.Write([]byte("+OK\r\n"))
			}
		case "get":
			value, ok := storage.Get(args[0].String())
			if ok {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}

	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"godis/resp"
	"godis/store"
	"strings"
)

var dataStore = store.NewStore()

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client connected: %s", conn.RemoteAddr())

	reader := resp.NewReader(conn)

	for {
		message, err := reader.ReadValue()
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", conn.RemoteAddr())
			} else {
				log.Printf("Error reading from connection %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		if message.Typ != resp.ARRAY || len(message.Array) == 0 {
			log.Printf("Invalid command from client %s", conn.RemoteAddr())
			conn.Write([]byte("-ERR invalid command\r\n"))
			continue
		}

		command := strings.ToUpper(message.Array[0].Bulk)
		args := message.Array[1:]

		switch command {
		case "SET":
			if len(args) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			key := args[0].Bulk
			value := []byte(args[1].Bulk)
			dataStore.Set(key, value)
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(args) != 1 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			key := args[0].Bulk
			value, ok := dataStore.Get(key)

			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				respStr := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				conn.Write([]byte(respStr))
			}
		
		case "LPUSH":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'lpush' command\r\n"))
				continue
			}
			key := args[0].Bulk
			values := make([][]byte, len(args)-1)
			for i := 1; i < len(args); i++ {
				values[i-1] = []byte(args[i].Bulk)
			}
			length := dataStore.LPush(key, values...)
			respStr := fmt.Sprintf(":%d\r\n", length)
			conn.Write([]byte(respStr))
		case "LPOP":
			if len(args) != 1 {
				conn.Write([]byte("-ERR wrong number of arguments for 'lpop' command\r\n"))
				continue
			}
			key := args[0].Bulk
			value, ok := dataStore.LPop(key)
			if !ok {
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			respStr := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
			conn.Write([]byte(respStr))
		case "PING":
			if len(args) > 1 {
				conn.Write([]byte("-ERR wrong number of arguments for 'ping' command\r\n"))
				continue
			}
			if len(args) == 0 {
				conn.Write([]byte("+PONG\r\n"))
			} else {
				pingMessage := args[0].Bulk
				respStr := fmt.Sprintf("+%s\r\n", pingMessage)
				conn.Write([]byte(respStr))
			}
		
		case "QUIT":
			if len(args) != 0 {
				conn.Write([]byte("-ERR wrong number of arguments for 'quit' command\r\n"))
				continue
			}
			conn.Write([]byte("+OK\r\n"))
			log.Printf("Client disconnected: %s", conn.RemoteAddr())
			return

		default:
			respStr := fmt.Sprintf("-ERR unknown command `%s`\r\n", command)
			conn.Write([]byte(respStr))
		}
	}
}

func main() {
	listenAddr := ":6379"

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Go-dis server started. Listening on %s", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

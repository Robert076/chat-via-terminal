package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]bool)
	clientsMu sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error when calling net.Listen in main()")
	}

	fmt.Println("Server is listening on port 8080...")
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Error accepting a new connection")
			break
		}

		fmt.Println("New connection accepted from ", conn)

		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from client")
			}
			break
		}
		message := string(buffer[:n])
		fmt.Println("Received: ", message)

		broadcastMessage(conn, message)
	}
}

func broadcastMessage(sender net.Conn, message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		if client != sender {
			_, _ = fmt.Fprintln(client, message)
		}
	}
}

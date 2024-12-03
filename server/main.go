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
	// ln is a pointer to net.TCPListener

	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()
	fmt.Println("Server started on port 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}

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
		// Read data from the connection
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break // Exit the loop if there's an error or EOF
		}

		// Extract the message (the first n bytes)
		message := string(buffer[:n])
		fmt.Println("Received:", message)

		// Broadcast the message to other clients (assuming a broadcast function)
		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if conn != sender {
			_, _ = fmt.Fprintln(conn, message)
		}
	}
}

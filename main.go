package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Client represents a connected user
type Client struct {
	conn net.Conn
	name string
}

// Variables to store clients and to control access
var clients = make(map[net.Conn]*Client)
var mutex = sync.Mutex()

// Broadcast sends a message to all connected clients
func Broadcast(message string, sender *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, client := range clients {
		if client != sender {
			client.conn.Write([]byte(message + "\n"))
		}
	}
}

// HandleClient manages incoming client connections and messages
func HandleClient(conn net.Conn) {
	defer conn.Close()

	// Prompt for username
	conn.Write([]byte("Enter your name: "))
	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)

	// Create a new client
	client := &Client{
		conn: conn,
		name: name,
	}

	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	// Notify other users
	Broadcast(fmt.Sprintf("%s has joined the chat", name), client)
	fmt.Printf("%s connected\n", name)

	// Read messages from the client and broadcast them
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			Broadcast(fmt.Sprintf("%s has left the chat", name), client)
			fmt.Printf("%s disconnected\n", name)
			return
		}

		// Broadcast the message
		Broadcast(fmt.Sprintf("%s: %s", name, strings.TrimSpace(message)), client)
	}
}

// StartServer initializes the chat server
func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go HandleClient(conn) // Handle each client in a new goroutine
	}
}

func main() {
	port := "8080" // Change to any port you want
	StartServer(port)
}

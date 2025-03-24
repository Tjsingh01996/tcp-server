package main

import (
	"fmt"
	"net"
	"time"

	"github.com/Tjsingh01996/tcp-server/https"
	"github.com/Tjsingh01996/tcp-server/utils"
)

var secretForChat = []byte("7aec2cb4d5c13ec40787a6286e103b42675fe090399e91916039ae687b566d97")

type ChatServer struct {
	connections map[string]connection
}

func NewChatServer() *ChatServer {

	return &ChatServer{
		connections: make(map[string]connection),
	}
}

func (chat *ChatServer) StartChatServer(port int) error {
	tcp := https.NewTcp()
	tcp.SetOnConnectNewConnection(chat.onConnectNewConnection)
	return tcp.Serve(port)
}

func (chat *ChatServer) onConnectNewConnection(conn net.Conn) error {
	payload := map[string]any{
		"time": time.Now().String(),
	}
	token, err := utils.GenerateJWT(secretForChat, payload, 0)
	if err != nil {
		return err
	}
	connection := connection{conn}
	go connection.ReadStream()
	chat.connections[token] = connection
	return nil
}

// individual connection //
// ====================================================== //

type connection struct {
	conn net.Conn
}

func (c *connection) ReadStream() {
	defer c.conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}

package main

import (
	"fmt"
	"net"
)

type Client struct {
	serverIp   string
	serverPort int
	name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	newClient := &Client{
		serverIp:   serverIp,
		serverPort: serverPort,
	}

	//connect server
	conn, err := net.Dial("tcp", fmt.Sprint("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.dial err:", err)
		return nil
	}
	newClient.conn = conn
	return newClient
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> server connect failed <<<<<")
		return
	}
	fmt.Println("connect success ! ")
	// start client mission
	select {}
}

package main

import (
	"fmt"
	"net"
)

type Client struct {
	Port int
	Conn net.Conn
}

func NewClient(port int) Client {
	return Client{
		Port: port,
	}
}

func (c *Client) Open() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", c.Port))
	if err != nil {
		panic(err)
	}
	c.Conn = conn
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Read() {
	buf := make([]byte, 1024*10)
	n, err := c.Conn.Read(buf)
	if err != nil {
		fmt.Println("ERR:", err.Error())
	} else {
		fmt.Println("READ:", string(buf), n, "bytes")
	}
}

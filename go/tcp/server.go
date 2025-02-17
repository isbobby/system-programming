package main

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	Port int
}

func NewServer(port int) Server {
	return Server{
		Port: port,
	}
}

func (s Server) Serve() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		conn.Write([]byte(fmt.Sprintf("Hello World %v", time.Now())))

		conn.Close()
	}
}

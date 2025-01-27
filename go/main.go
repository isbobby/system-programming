package main

import (
	"go-sockets/tcp"
	"sync"
)

func main() {
	tcpServer := tcp.NewServer(9000)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		tcpServer.Serve()
	}()

	tcpClient := tcp.NewClient(9000)
	tcpClient.Open()
	defer tcpClient.Close()
	for i := 0; i < 10; i++ {
		tcpClient.Read()
	}
	wg.Wait()
}

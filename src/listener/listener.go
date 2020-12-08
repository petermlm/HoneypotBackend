package listener

import (
	"context"
	"fmt"
	"log"
	"net"
)

func Start(ctx context.Context, ports []string) {
	waitChans := make([]chan bool, len(ports))

	for i, port := range ports {
		waitChans[i] = make(chan bool)
		go listen(ctx, port, waitChans[i])
	}

	for _, ch := range waitChans {
		<-ch
	}
}

func listen(ctx context.Context, port string, wait chan bool) {
	defer func() { wait <- true }()

	listener, err := createListener(port)
	if err != nil {
		log.Printf("Can't create listener for port %s\n", port)
		return
	}

	acceptChan := make(chan net.Conn)
	acceptFunc := func() {
		conn, err := listener.Accept()
		if err != nil {
			acceptChan <- nil
		} else {
			acceptChan <- conn
		}
	}

	log.Printf("Listening on %s\n", port)
	for {
		go acceptFunc()
		select {
		case conn := <-acceptChan:
			if conn == nil {
				continue
			}
			log.Println("Conn", port)
			conn.Close()
		case <-ctx.Done():
			return
		}
	}
}

func createListener(port string) (net.Listener, error) {
	connAddr := ":" + port
	l, err := net.Listen("tcp", connAddr)

	if err != nil {
		return nil, fmt.Errorf("Error listening: %w", err.Error())
	}

	return l, nil
}

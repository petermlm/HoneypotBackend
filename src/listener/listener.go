package listener

import (
	"context"
	"fmt"
	"honeypot/timelines"
	"log"
	"net"
	"time"
)

func Start(ctx context.Context, ports []string) {
	waitChans := make([]chan bool, len(ports))

	tl := timelines.NewTimelinesWriter()

	for i, port := range ports {
		waitChans[i] = make(chan bool)
		go listen(ctx, waitChans[i], tl, port)
	}

	for _, ch := range waitChans {
		<-ch
	}
}

func listen(ctx context.Context, wait chan bool, tl timelines.TimelinesWriter, port string) {
	defer func() { wait <- true }()

	listener, err := createListener(port)
	if err != nil {
		log.Printf("Can't create listener for port %s\n", port)
		return
	}
	defer listener.Close()

	acceptChan := make(chan net.Conn)
	acceptFunc := func() {
		conn, err := listener.Accept()
		if err == nil && conn != nil {
			acceptChan <- conn
		}
	}

	log.Printf("Listening on %s\n", port)
	for {
		go acceptFunc()
		select {
		case conn := <-acceptChan:
			registerConnAttemp(tl, conn, port)
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

func registerConnAttemp(tl timelines.TimelinesWriter, conn net.Conn, port string) {
	log.Println("Conn", port)

	addr := conn.RemoteAddr().String()
	point := &timelines.ConnAttemp{
		Time: time.Now(),
		Port: port,
		Addr: addr,
	}
	tl.InsertConnAttemp(point)
	conn.Close()
}

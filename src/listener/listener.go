package listener

import (
	"context"
	"fmt"
	"honeypot/timelines"
	"log"
	"net"
	"time"
)

func Start(ctx context.Context, ports []string) error {
	lenPorts := len(ports)
	waitChans := make([]chan bool, lenPorts)
	cancelChans := make([]chan bool, lenPorts)

	tl := timelines.NewTimelinesWriter()
	errorsCh := tl.Errors()

	for i, port := range ports {
		waitChans[i] = make(chan bool, 1)
		cancelChans[i] = make(chan bool, 1)
		go listen(ctx, waitChans[i], cancelChans[i], tl, port)
	}

	var err error
	select {
	case e := <-errorsCh:
		err = e
	case <-ctx.Done():
	}

	for _, ch := range cancelChans {
		ch <- true
	}
	for _, ch := range waitChans {
		<-ch
	}
	return err
}

func listen(ctx context.Context, wait chan bool, cancel chan bool, tl timelines.TimelinesWriter, port string) {
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
		case <-cancel:
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
	connAttemp, err := timelines.NewConnAttemp(time.Now(), port, addr)
	if err != nil {
		// TODO: Handle it
	}

	tl.InsertConnAttemp(connAttemp)
	conn.Close()
}

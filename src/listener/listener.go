package listener

import (
	"context"
	"fmt"
	"honeypot/queue"
	"honeypot/settings"
	"honeypot/timelines"
	"log"
	"net"
	"time"
)

func Start(ctx context.Context, ports []string) error {
	var err error

	lenPorts := len(ports)
	waitChans := make([]chan bool, lenPorts)
	cancelChans := make([]chan bool, lenPorts)

	publisher, err := queue.NewPublisher(settings.RabbitmqTaskProcessConnAttemp)
	if err != nil {
		return err
	}

	tl := timelines.NewTimelinesWriter()
	errorsCh := tl.Errors()

	for i, port := range ports {
		waitChans[i] = make(chan bool, 1)
		cancelChans[i] = make(chan bool, 1)
		go listen(ctx, waitChans[i], cancelChans[i], tl, publisher, port)
	}

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

func listen(ctx context.Context, wait chan bool, cancel chan bool, tl timelines.TimelinesWriter, publisher queue.Publisher, port string) {
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
			sendToProcessor(publisher, conn, port)
		case <-cancel:
			return
		}
	}
}

func createListener(port string) (net.Listener, error) {
	connAddr := ":" + port
	l, err := net.Listen("tcp", connAddr)

	if err != nil {
		return nil, fmt.Errorf("Error listening: %v", err.Error())
	}

	return l, nil
}

func sendToProcessor(publisher queue.Publisher, conn net.Conn, port string) {
	defer conn.Close()

	// Create conn attemp
	addr := conn.RemoteAddr().String()
	connAttemp, err := timelines.NewConnAttemp(time.Now(), port, addr)
	if err != nil {
		return
	}
	defer publisher.Publish(connAttemp)

	// Try to read something
	deadline := time.Now().Add(time.Second * 5)
	conn.SetReadDeadline(deadline)

	b := make([]byte, 1024*4)
	n, err := conn.Read(b)
	if err != nil {
		return
	}

	if n >= 0 {
		connAttemp.Bytes = b
	}
}

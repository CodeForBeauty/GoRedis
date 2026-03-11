package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

func StartServer(port int) error {
	fmt.Printf("Starting server on port: %d...\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	var toProcess = make(chan struct {
		message    string
		connection net.Conn
	}, 5)

	if err != nil {
		return err
	}

	go RunProcessor(toProcess)

	defer listener.Close()

	for {
		connection, err := listener.Accept()

		if err != nil {
			return err
		}

		go HandleConnection(connection, toProcess)
	}
}

func HandleConnection(conn net.Conn, toProcess chan struct {
	message    string
	connection net.Conn
}) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	if reader == nil {
		return
	}

	for {
		message, err := reader.ReadString('\n')

		if err != nil && !errors.Is(err, io.EOF) || message == "" {
			break
		}

		toProcess <- struct {
			message    string
			connection net.Conn
		}{message: message, connection: conn}
	}
}

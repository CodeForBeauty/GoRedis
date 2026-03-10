package server

import (
	"bufio"
	"fmt"
	"net"
)

func StartServer(port int) error {
	fmt.Printf("Starting server on port: %d...\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()

		if err != nil {
			return err
		}

		go HandleConnection(connection)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	if reader == nil {
		return
	}

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		fmt.Printf("Message: %s\n", message)
	}
}

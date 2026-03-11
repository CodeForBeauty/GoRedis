package internal

import (
	"fmt"
	"net"
)

func RunProcessor(toProcess chan struct {
	message    string
	connection net.Conn
}) {
	dbServer := MakeServer()

	for {
		command := <-toProcess

		fmt.Printf("Processing: %s\n", command.message)
		output, err := dbServer.ProcessCommand(command.message)
		if err != nil {
			command.connection.Write([]byte("FAILED"))
		} else if output != "" {
			command.connection.Write([]byte(output))
		} else {
			command.connection.Write([]byte("OK"))
		}
	}
}

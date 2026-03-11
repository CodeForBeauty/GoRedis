package internal

import (
	"fmt"
	"net"
)

var toProcess = make(chan struct {
	message    string
	connection net.Conn
}, 5)

func RunProcessor() {
	for {
		command := <-toProcess

		fmt.Printf("Processing: %s\n", command.message)
		output, err := ProcessCommand(command.message)
		if err == nil && output != "" {
			command.connection.Write([]byte(output))
		}
	}
}

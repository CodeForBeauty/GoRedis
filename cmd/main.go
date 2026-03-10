package main

import (
	server "github.com/CodeForBeauty/GoRedis/internal"
)

func main() {
	server.StartServer(8080)
}

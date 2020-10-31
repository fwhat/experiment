package main

import "auto_proxy/server"

func main() {
	panic(server.NewServer().Serve(":4400"))
}

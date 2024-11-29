package main

import server "AuthServ/Server"

func init() {
	server.InitServer()
}

func main() {
	server.StartServer()
}

package main

import "github.com/wellingtonlope/calculator-api/cmd/http/server"

func main() {
	server := server.New()
	server.Logger.Fatal(server.Start(":8080"))
}

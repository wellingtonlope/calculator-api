package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/wellingtonlope/calculator-api/cmd/http/server"
)

func main() {
	server := server.New()
	server.Logger.Fatal(server.Start(":8080"))
}

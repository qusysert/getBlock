package main

import (
	"log"

	"getBlock/internal/server"

	"getBlock/config"
)

func main() {
	cfg := config.LoadConfig()

	srv := server.NewServer(&cfg)
	log.Println("Server is running on port 8080")
	srv.Run()
}

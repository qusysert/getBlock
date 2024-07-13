package main

import (
	"getBlock/config"
	"getBlock/internal/server"
)

func main() {
	cfg := config.LoadConfig()

	srv := server.NewServer(&cfg)
	srv.Run()
}

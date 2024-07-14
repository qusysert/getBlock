package main

import (
	"getBlock/config"
	"getBlock/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	srv := server.NewServer(&cfg)
	srv.Run()
}

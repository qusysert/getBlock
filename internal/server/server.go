package server

import (
	"log"
	"net/http"

	"getBlock/config"
	"getBlock/internal/getblock"
)

type Server struct {
	client *getblock.Client
}

func NewServer(cfg *config.Config) *Server {
	client := getblock.NewClient(cfg.ApiKey)
	return &Server{client: client}
}

func (s *Server) Run() {
	http.HandleFunc("/balance", s.balanceHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

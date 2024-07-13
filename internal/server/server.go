package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"getBlock/config"
	"getBlock/internal/balance"
	"getBlock/internal/getblock"
)

type Server struct {
	client *getblock.Client
}

func NewServer(cfg *config.Config) *Server {
	client := getblock.NewClient(cfg.ApiKey)
	return &Server{client: client}
}

func (s *Server) balanceHandler(w http.ResponseWriter, r *http.Request) {
	latestBlockNumberHex, err := s.client.GetLatestBlockNumber()
	if err != nil {
		http.Error(w, "Failed to get latest block number", http.StatusInternalServerError)
		return
	}

	// Convert hex string to int
	latestBlockNumber, err := strconv.ParseInt(latestBlockNumberHex[2:], 16, 64)
	if err != nil {
		http.Error(w, "Failed to parse latest block number", http.StatusInternalServerError)
		return
	}

	var blocks []getblock.Block
	for i := int64(0); i < 100; i++ {
		blockNumber := fmt.Sprintf("0x%x", latestBlockNumber-i)
		block, err := s.client.GetBlockByNumber(blockNumber)
		if err != nil {
			http.Error(w, "Failed to get block data", http.StatusInternalServerError)
			return
		}
		blocks = append(blocks, block)
	}

	balanceChanges := balance.CalculateBalanceChanges(blocks)
	sort.Slice(balanceChanges, func(i, j int) bool {
		return balanceChanges[i].Change.Cmp(balanceChanges[j].Change) > 0
	})

	response, err := json.Marshal(balanceChanges[0])
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *Server) Run() {
	http.HandleFunc("/balance", s.balanceHandler)
	http.ListenAndServe(":8080", nil)
}

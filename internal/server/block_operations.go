package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"getBlock/internal/balance"
	"getBlock/internal/getblock"
)

func (s *Server) fetchBlock(blockNumber int64, wg *sync.WaitGroup, ch chan<- getblock.Block, errCh chan<- error, rateLimiter <-chan time.Time) {
	defer wg.Done()
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	retries := 0

	for {
		<-rateLimiter
		block, err := s.client.GetBlockByNumber(blockNumberHex)
		if err == nil {
			ch <- block
			return
		}
		if retries >= 5 {
			errCh <- fmt.Errorf("failed to fetch block %s: %w", blockNumberHex, err)
			return
		}
		retries++
		time.Sleep(time.Duration(retries) * time.Second)
	}
}

func (s *Server) balanceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling balance request")

	latestBlockNumberHex, err := s.client.GetLatestBlockNumber()
	if err != nil {
		log.Printf("Failed to get latest block number: %v\n", err)
		http.Error(w, "Failed to get latest block number", http.StatusInternalServerError)
		return
	}

	log.Printf("Latest block number (hex): %s\n", latestBlockNumberHex)

	latestBlockNumber, err := strconv.ParseInt(latestBlockNumberHex[2:], 16, 64)
	if err != nil {
		log.Printf("Failed to parse latest block number: %v\n", err)
		http.Error(w, "Failed to parse latest block number", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	blockCh := make(chan getblock.Block, 100)
	errCh := make(chan error, 100)
	rateLimiter := time.Tick(time.Second / 60)

	for i := int64(0); i < 100; i++ {
		wg.Add(1)
		go s.fetchBlock(latestBlockNumber-i, &wg, blockCh, errCh, rateLimiter)
	}

	go func() {
		wg.Wait()
		close(blockCh)
		close(errCh)
	}()

	var blocks []getblock.Block
	for block := range blockCh {
		blocks = append(blocks, block)
	}

	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("Error fetching block data: %v\n", err)
			http.Error(w, "Failed to get block data", http.StatusInternalServerError)
			return
		}
	default:
	}

	balanceChanges := balance.CalculateBalanceChanges(blocks)
	sort.Slice(balanceChanges, func(i, j int) bool {
		return balanceChanges[i].Change.Cmp(balanceChanges[j].Change) > 0
	})

	response, err := json.Marshal(balanceChanges[0])
	if err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

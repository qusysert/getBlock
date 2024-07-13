package balance

import (
	"math/big"

	"getBlock/internal/getblock"
)

type BalanceChange struct {
	Address string
	Change  *big.Int
}

func CalculateBalanceChanges(blocks []getblock.Block) []BalanceChange {
	balanceChanges := make(map[string]*big.Int)

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			value, _ := new(big.Int).SetString(tx.Value[2:], 16)
			gasPrice, _ := new(big.Int).SetString(tx.GasPrice[2:], 16)

			if _, exists := balanceChanges[tx.From]; !exists {
				balanceChanges[tx.From] = big.NewInt(0)
			}
			balanceChanges[tx.From].Sub(balanceChanges[tx.From], value)
			balanceChanges[tx.From].Sub(balanceChanges[tx.From], gasPrice)

			if _, exists := balanceChanges[tx.To]; !exists {
				balanceChanges[tx.To] = big.NewInt(0)
			}
			balanceChanges[tx.To].Add(balanceChanges[tx.To], value)
		}
	}

	changes := make([]BalanceChange, 0, len(balanceChanges))
	for addr, change := range balanceChanges {
		changes = append(changes, BalanceChange{
			Address: addr,
			Change:  change,
		})
	}

	return changes
}

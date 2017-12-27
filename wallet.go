package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/grrrben/golog"
	"math/big"
	"time"
)

type wallet struct {
	hash   string
	credit float64
}

// createWallet creates a wallet with a hash and 0 credits
func createWallet() wallet {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, time.Now().Unix())
	if err != nil {
		golog.Warningf("Could not createWallet. Msg: %s", err)
	}

	w := wallet{
		hash:   fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())),
		credit: 0,
	}
	return w
}

// getWalletCredits Loops all blocks/transactions and checks for the credits that are send or received.
// Also loops the current pending transactions that are not mined yet. Of _this_ client...
// returns the total amount of credits that are currently in the wallet
func getWalletCredits(hash string) float64 {
	var sum big.Float
	for _, block := range bc.Chain {
		for _, transaction := range block.Transactions {
			if transaction.Recipient == hash {
				sum.Add(&sum, big.NewFloat(transaction.Amount))
			}
			if transaction.Sender == hash {
				sum.Sub(&sum, big.NewFloat(transaction.Amount))
			}
		}
	}

	for _, pendingTransactions := range bc.Transactions {
		if pendingTransactions.Recipient == hash {
			sum.Add(&sum, big.NewFloat(pendingTransactions.Amount))
		}
		if pendingTransactions.Sender == hash {
			sum.Sub(&sum, big.NewFloat(pendingTransactions.Amount))
		}
	}

	credits, _ := sum.Float64()
	return credits
}

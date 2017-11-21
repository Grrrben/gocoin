package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/grrrben/golog"
	"time"
)

type wallet struct {
	hash   string
	credit float64
}

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
func getWalletCredits(hash string) float32 {
	var credit float32
	for _, block := range bc.Chain {
		for _, transaction := range block.Transactions {
			if transaction.Recipient == hash {
				credit = credit + transaction.Amount
			}
			if transaction.Sender == hash {
				credit = credit - transaction.Amount
			}
		}
	}

	for _, pendingTransactions := range bc.Transactions {
		if pendingTransactions.Recipient == hash {
			credit = credit + pendingTransactions.Amount
		}
		if pendingTransactions.Sender == hash {
			credit = credit - pendingTransactions.Amount
		}
	}
	return credit
}

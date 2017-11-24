package main

import (
	"errors"
	"log"
	"regexp"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float32
}

// validHash checks a hash for length and regex of the hex.
// It does _not_ check for existince of a wallet with this specific hash.
func validHash(hash string) bool {
	// fad5e7a92f1c43b1523614336a07f98b894bb80fee06b6763b50ab03b597d5f4
	regex, err := regexp.Compile(`[a-f0-9]{64}`)

	if err != nil {
		log.Fatal("Could not compile regex")
	}
	if regex.MatchString(hash) {
		return true
	} else {
		return false
	}
}

// checkTransaction performs multiple checks on a transaction
func checkTransaction(tr Transaction) (success bool, err error) {
	if !validHash(tr.Sender) {
		return false, errors.New("invalid transaction (sender invalid)")
	} else if !validHash(tr.Recipient) {
		return false, errors.New("invalid transaction (recipient invalid)")
	} else if getWalletCredits(tr.Sender) < tr.Amount {
		return false, errors.New("invalid transaction (insufficient credit)")
	}
	return true, nil
}

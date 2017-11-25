package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float32
	Time      int64
}

type hashable interface {
	getHash() string
}

func (tr Transaction) getHash() string {
	str := tr.Sender + tr.Recipient + fmt.Sprintf("%.8f", tr.Amount) + fmt.Sprintf("%d", tr.Time)
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))

}

// checkHashes checks if the hashes of 2 objects are the same
// objects should have interface hashable.
func checkHashes(first hashable, second hashable) bool {
	if first.getHash() == second.getHash() {
		return true
	}
	return false
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

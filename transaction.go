package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grrrben/golog"
	"log"
	"net/http"
	"regexp"
)

type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Message   string  `json:"message"`
	Time      int64   `json:"time"`
}

type hashable interface {
	getHash() string
}

// getHash a unique hash for a transaction
func (tr Transaction) getHash() string {
	str := tr.Sender + tr.Recipient + fmt.Sprintf("%.8f", tr.Amount) + fmt.Sprintf("%d", tr.Time)
	sha := sha256.New()
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil))

}

// checkHashesEqual checks if the hashes of 2 objects are the same
// objects should have interface hashable.
func checkHashesEqual(first hashable, second hashable) bool {
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
	} else if tr.Sender != zerohash && getWalletCredits(tr.Sender) < tr.Amount {
		return false, errors.New("invalid transaction (insufficient credit)")
	}
	return true, nil
}

// announceTransaction distributes new transaction in the network
// It is preferably done in a goroutine.
func announceTransaction(cl Client, tr Transaction) {
	defer golog.Flush()
	url := fmt.Sprintf("%s/transaction/distributed", cl.getAddress())

	transactionAndSender := map[string]interface{}{"transaction": tr, "sender": me.getAddress()}
	golog.Infof("transactionAndSender to be distributed:\n %v", transactionAndSender)
	payload, err := json.Marshal(transactionAndSender)
	if err != nil {
		golog.Errorf("Could not marshall transaction or client. Msg: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		golog.Warningf("Request setup error: %s", err)
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Warningf("POST request error: %s", err)
		// I don't want to panic here, but it might be a good idea to
		// remove the client from the list
	} else {
		defer resp.Body.Close()
		golog.Info("Transaction distributed")
	}
}

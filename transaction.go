package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"encoding/json"
	"github.com/grrrben/golog"
	"net/http"
	"bytes"
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

// getHash a unique hash for a transaction
func (tr Transaction) getHash() string {
	str := tr.Sender + tr.Recipient + fmt.Sprintf("%.8f", tr.Amount) + fmt.Sprintf("%d", tr.Time)
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))

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
	} else if getWalletCredits(tr.Sender) < tr.Amount {
		return false, errors.New("invalid transaction (insufficient credit)")
	}
	return true, nil
}

// announceTransaction distributes new transaction in the network
// It is preferably done in a goroutine.
func announceTransaction(cl Client, tr Transaction) {
	url := fmt.Sprintf("%s/transaction/distributed", cls.getAddress(cl))

	transactionAndSender := map[string]interface{}{"transaction": tr, "sender": cls.getAddress(me)}
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
	}
	defer resp.Body.Close()
}

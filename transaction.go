package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/grrrben/glog"
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
	str := fmt.Sprintf("%s%s%.8f%d", tr.Sender, tr.Recipient, tr.Amount, tr.Time)
	sha := sha256.New()
	sha.Write([]byte(str))
	return fmt.Sprintf("%x", sha.Sum(nil))

}

// checkHashesEqual checks if the hashes of 2 objects are the same
// objects should have interface hashable.
func checkHashesEqual(first, second hashable) bool {
	return first.getHash() == second.getHash()
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
func announceTransaction(node Node, tr Transaction) {
	defer glog.Flush()
	url := fmt.Sprintf("%s/transaction/distributed", node.getAddress())

	transactionAndSender := make(map[string]interface{}, 2)
	transactionAndSender["transaction"] = tr
	transactionAndSender["sender"] = me.getAddress()
	glog.Infof("transactionAndSender to be distributed:\n %v", transactionAndSender)
	payload, err := json.Marshal(transactionAndSender)
	if err != nil {
		glog.Errorf("Could not marshall transaction or node. Msg: %s", err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		glog.Errorf("Request setup error: %s", err.Error())
	} else {
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			glog.Warningf("POST request error: %s", err.Error())
			// I don't want to panic here, but it might be a good idea to
			// remove the node from the list (todo)
		} else {
			defer resp.Body.Close()
			glog.Info("Transaction distributed")
		}
	}
}

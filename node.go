package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grrrben/glog"
)

// this is me, a node
var me Node

type Node struct {
	Hostname string `json:"hostname"`
	Protocol string `json:"protocol"`
	Port     uint16 `json:"port"`
	Name     string `json:"name"`
	Hash     string `json:"hash"`
}

// greet makes a call to a node cl to make this node known within the network.
func greet(cl Node) {
	// POST to /node
	url := fmt.Sprintf("%s/node", cl.getAddress())
	payload, err := json.Marshal(me)
	if err != nil {
		glog.Warning("Could not marshall node: Me")
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		glog.Warningf("Request setup error: %s", err)
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Warningf("POST request error: %s", err)
		// I don't want to panic here, but it might be a good idea to
		// remove the node from the list
	} else {
		resp.Body.Close()
	}
}

// createWallet Creates a wallet and sets the hash of the new wallet on the Node.
// Is is done only once. As soon as the wallet hash is set this function does nothing.
// If a nodes mines a block, the incentive is sent to this wallet address
func (cl *Node) createWallet() {
	if !hasValidHash(cl) {
		wallet := createWallet()
		cl.Hash = wallet.hash
	}
}

// to make each Node a Hashable (interface)
func (cl Node) getHash() string {
	return cl.Hash
}

// getAddress returns (URI) address of a node.
func (cl Node) getAddress() string {
	return fmt.Sprintf("%s%s:%d", cl.Protocol, cl.Hostname, cl.Port)
}

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

// greet makes a call to a node to make this node known within the network.
func greet(node Node) {
	url := fmt.Sprintf("%s/node", node.getAddress())
	payload, err := json.Marshal(me)
	if err != nil {
		glog.Panicf("Could not marshall node: Me; %s", err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		glog.Panicf("Request setup error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Warningf("POST request error: %s", err)
		// I don't want to panic here, but it might be a good idea to
		// remove the node from the list (todo)
	} else {
		resp.Body.Close()
	}
}

// createWallet Creates a wallet and sets the hash of the new wallet on the Node.
// Is is done only once. As soon as the wallet hash is set this function does nothing.
// If a nodes mines a block, the incentive is sent to this wallet address
func (node *Node) createWallet() {
	if !hasValidHash(node) {
		wallet := createWallet()
		node.Hash = wallet.hash
	}
}

// getHash; to make each Node a Hashable (interface)
func (node Node) getHash() string {
	return node.Hash
}

// getAddress returns (URI) address of a node.
func (node Node) getAddress() string {
	return fmt.Sprintf("%s%s:%d", node.Protocol, node.Hostname, node.Port)
}

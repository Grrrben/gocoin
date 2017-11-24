package main

import (
	"flag"
	"testing"
)

// Testing a couple of blockchain related functions on a base chain
// with just the genesis block

func init() {
	// setup
	flag.Parse()
	cls = initClients()
	me = Client{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8000,
		Name:     "client1",
		Hash:     createClientHash(),
	}
	cls.addClient(me)
}

func TestInitBlockchain(t *testing.T) {

	bc = initBlockchain()
	if len(bc.Chain) != 1 {
		t.Errorf("Chainlength incorrect got: %d, want: %d.", len(bc.Chain), 1)
	}

	if bc.Chain[len(bc.Chain)-1].Index != 1 {
		t.Errorf("Index of genesis block incorrect got: %d, want: %d.", bc.Chain[len(bc.Chain)-1].Index, 1)
	}

	if len(bc.Transactions) != 0 {
		t.Errorf("New blockchain should not have transactions. Got: %d transactions.", len(bc.Transactions))
	}
}

func TestNewTransaction(t *testing.T) {
	transaction := Transaction{
		"sender",
		"receiver",
		1,
	}

	newBlockIndex := bc.newTransaction(transaction)
	lastBlockIndex := bc.lastBlock().Index + 1

	if newBlockIndex != lastBlockIndex {
		t.Errorf("Index of new block fails. Got %d, Want %d", newBlockIndex, lastBlockIndex)
	}
}

// TestValidate Current blockchain just has the Genesis block
// should always be valid
func TestValidate(t *testing.T) {
	valid := bc.validate()
	if !valid {
		t.Error("Blockchain is invalid.")
	}
}

func TestLastBlock(t *testing.T) {
	block := bc.lastBlock()
	if block.Index != 1 {
		t.Errorf("Last block index should be 2, got %d.", block.Index)
	}
}

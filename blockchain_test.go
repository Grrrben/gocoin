package main

import (
	"flag"
	"testing"
	"time"

	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/grrrben/glog"
)

// Testing a couple of blockchain related functions on a base chain
// with just the genesis block

func init() {
	// setup
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Could not set a logdir. Msg %s", err)
	}

	glog.SetLogFile(fmt.Sprintf("%s/log/blockchain.log", dir))
	glog.SetLogLevel(glog.Log_level_error)

	nodes = initNodes()
	me = Node{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8000,
		Name:     "node1",
	}
	me.createWallet()
	nodes.addNode(&me)
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
		"",
		time.Now().UnixNano(),
	}

	_, e := bc.newTransaction(transaction)

	if e == nil {
		t.Error("Expected an error when adding an invalid transaction")
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

package main

import (
	"testing"
	"flag"
)

func init() {
	// setup
	flag.Parse()
	cls = initClients()
	me = Client{
		Protocol: "http://",
		Ip:       "127.0.0.1",
		Port:     8000,
		Name:     "client1",
		Hash:     createClientHash("127.0.0.1", 8000, "test client"),
	}
	cls.addClient(me)
}

func TestInitBlockchain(t *testing.T) {

	bc = initBlockchain()
	if len(bc.Chain) != 1 {
		t.Errorf("Chainlength incorrect got: %d, want: %d.", len(bc.Chain), 1)
	}

	if bc.Chain[len(bc.Chain) - 1].Index != 1 {
		t.Errorf("Index of genesis block incorrect got: %d, want: %d.", bc.Chain[len(bc.Chain) - 1].Index, 1)
	}

	if len(bc.Transactions) != 0 {
		t.Errorf("New blockchain should not have transactions. Got: %d transactions.", len(bc.Transactions))
	}
}

package main

import (
	"testing"
)

var node Node

func init() {
	node = Node{
		Protocol: "http://",
		Hostname: "localhost",
		Port:     8001,
		Name:     "node_x",
	}
	node.createWallet()
}

func TestCreateNodeWallet(t *testing.T) {
	if !hasValidHash(node) {
		t.Error("Node has no hash of wallet.")
	}
}

func TestGetNodeHash(t *testing.T) {
	if node.Hash != node.getHash() {
		t.Errorf("Incorrect node Hash: Expected %s, got %s.", node.Hash, node.getHash())
	}
}

func TestGetNodeAddress(t *testing.T) {
	if node.getAddress() != "http://localhost:8001" {
		t.Errorf("Incorrect node address: Expected http://localhost:8001, got %s.", node.getAddress())
	}
}

func TestGreet(t *testing.T) {
	me = Node{
		Protocol: "http://",
		Hostname: "localhost",
		Port:     8000,
		Name:     "this is me",
	}
	me.createWallet()
	greet(node)
}

package main

import (
	"testing"
)

var client Client

func init() {
	client = Client{
		Protocol: "http://",
		Hostname: "localhost",
		Port:     8001,
		Name:     "client1",
	}
	client.Hash = client.createWallet()
}

func TestCreateClientWallet(t *testing.T) {
	if !hasValidHash(client) {
		t.Error("Client has no hash of wallet.")
	}
}

func TestGetClientHash(t *testing.T) {
	if client.Hash != client.getHash() {
		t.Errorf("Incorrect client Hash: Expected %s, got %s.", client.Hash, client.getHash())
	}
}

func TestGetClientAddress(t *testing.T) {
	if client.getAddress() != "http://localhost:8000" {
		t.Errorf("Incorrect client address: Expected http://localhost:8000, got %s.", client.getAddress())
	}
}

func TestGreet(t *testing.T) {
	me = Client{
		Protocol: "http://",
		Hostname: "localhost",
		Port:     8000,
		Name:     "this is me",
	}
	me.Hash = me.createWallet()
	greet(client)
}

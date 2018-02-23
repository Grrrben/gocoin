package main

import (
	"testing"
)

func TestSortMapDescending(t *testing.T) {

	list := make(map[interface{}]int, 4)
	list["a"] = 1
	list["b"] = 2
	list["c"] = 3
	list["d"] = 4

	sorted := sortMapDescending(list)
	first := sorted[0]

	if first.Key != "d" {
		t.Errorf("SortMapDescending test failed. Expected 'd', got %s", first.Key)
	}
}

func TestSortMapAscending(t *testing.T) {

	list := make(map[interface{}]int, 4)
	list["a"] = 8
	list["b"] = 7
	list["c"] = 6
	list["d"] = 5

	sorted := sortMapAscending(list)
	first := sorted[0]

	if first.Key != "d" {
		t.Errorf("SortMapAscending test failed. Expected 'd', got %s", first.Key)
	}
}

func TestHasValidHash(t *testing.T) {
	testClient := Client{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8000,
		Name:     "client1",
	}
	testClient.createWallet()
	if !hasValidHash(testClient) {
		t.Errorf("Hash of client invalid (%s).", testClient.Hash)
	}
}

func TestValidHash(t *testing.T) {
	hash := "fad5e7a92f1c43b1523614336a07f98b894bb80fee06b6763b50ab03b597d5f4"
	if !validHash(hash) {
		t.Errorf("Valid hash '%s' was not tested positive", hash)
	}

	invalidHash := "notagoodhash"
	if validHash(invalidHash) {
		t.Errorf("Invalid hash '%s' was tested positive", invalidHash)
	}
}

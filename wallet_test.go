package main

import (
	"testing"
	"regexp"
)

func TestCreateWallet(t *testing.T) {
	wallet := createWallet()

	regex, _ := regexp.Compile(`[a-f0-9]{64}`)
	if !regex.MatchString(wallet.hash) {
		t.Error("Wallet hash not according to '[a-f0-9]{64}'.")
	}

	if wallet.credit != 0 {
		t.Fail()
	}
}

package main

import (
	"strings"
	"testing"
)

func TestGetHash(t *testing.T) {

	transaction := Transaction{
		"sender",
		"recipient",
		1.2,
		"message",
		0,
	}

	hash := transaction.getHash()

	if hash != "0c77fd2a8a8f74f6f547144b074f5b8fec39dd686209acaec2d70587e8f28b8a" {
		t.Errorf("transaction.getHash() test failed. Expected '0c77fd2a8a8f74f6f547144b074f5b8fec39dd686209acaec2d70587e8f28b8a', got %s", hash)
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

func TestCheckTransaction(t *testing.T) {
	// an invalid transaction
	tr := Transaction{
		"sender",
		"recipient",
		1.2,
		"message",
		0,
	}

	successSenderInvalid, errSenderInvalid := checkTransaction(tr)

	if successSenderInvalid {
		t.Error("checkTransaction: Invalid sender in transaction should result in false.")
	}

	if !strings.Contains(errSenderInvalid.Error(), "sender invalid") {
		t.Errorf("Expected error 'sender invalid', got %s.", errSenderInvalid.Error())
	}

	// setting a valid sender
	tr.Sender = "fad5e7a92f1c43b1523614336a07f98b894bb80fee06b6763b50ab03b597d5f4"

	successRecipientInvalid, errRecipientInvalid := checkTransaction(tr)

	if successRecipientInvalid {
		t.Error("checkTransaction: Invalid recipient in transaction should result in false.")
	}

	if !strings.Contains(errRecipientInvalid.Error(), "recipient invalid") {
		t.Errorf("Expected error 'recipient invalid', got %s.", errRecipientInvalid.Error())
	}
}

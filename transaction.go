package main

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float32
}

func validHash(hash string) bool {
	// fad5e7a92f1c43b1523614336a07f98b894bb80fee06b6763b50ab03b597d5f4
	if len(hash) == 64 {
		// todo check regex [a-f0-9]{64}
		return true
	} else {
		return false
	}

}

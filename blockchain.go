package main

import (
	"time"
	"fmt"
	"bytes"
	"encoding/binary"
	"crypto/sha256"
	"encoding/json"
)

// how many 0's do we want to check
const hashDifficulty int8 = 4

type Blockchain struct {
	Chain        []Block
	Transactions []Transaction
	Nodes        []string
}

type chainService interface {
	newBlock() bool
	newTransaction() bool
	hash(Block) string
	lastBlock() Block
	proofOfWork(lastProof int64) int64
	validProof(proof int64, lastProof int64) bool
	validate() bool
}

// newTransaction will create a Transaction to go into the next Block to be mined.
// The Transaction is stored in the Blockchain obj.
// Returns (int) the Index of the Block that will hold this Transaction
func (bc *Blockchain) newTransaction(tr Transaction) int64 {
	bc.Transactions = append(bc.Transactions, tr)
	fmt.Println("Transaction added")
	return bc.lastBlock().Index + 1
}

// hash Creates a SHA-256 hash of a Block
func hash(b Block) string {
	fmt.Printf("hashing block %d\n", b.Index)

	// Data for binary.Write must be a fixed-size value or a slice of fixed-size values,
	// or a pointer to such data.
	// @todo Marshalling the struct to json is a workaround... But it works
	// @todo might be able to fix it with a char(length) instead of string?
	jsonblock, errr := json.Marshal(b)
	if errr != nil {
		fmt.Printf("Error: %s", errr)
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, jsonblock)
	if err != nil {
		fmt.Println("Could not compute hash")
		fmt.Println(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())) // %x; base 16, with lower-case letters for a-f
}

// lastBlock returns the last Block in the Chain
func (bc *Blockchain) lastBlock() Block {
	return bc.Chain[len(bc.Chain) -1]
}

func (bc *Blockchain) proofOfWork(lastProof int64) int64 {
	// Simple Proof of Work Algorithm:
	// - Find a number p' such that hash(lp') contains leading 4 zeroes, where
    // - l is the previous Proof, and p' is the new Proof
	var proof int64 = 0
	i := 0;
	for !bc.validProof(lastProof, proof) {
		proof += 1
		i++;
	}
	fmt.Printf("Proof found in %d cycles (difficulty %d)\n", i, hashDifficulty)
	return proof

}

// validProof is called until it finds an acceptable hash and returns true
func (bc *Blockchain) validProof(proof int64, lastProof int64) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := fmt.Sprintf("%x", sha256.Sum256([]byte(guess)))

	var i int8
	hashString := ""
	for i = 0; i < hashDifficulty; i++ {
		hashString = hashString + "0"
	}

	if guessHash[:hashDifficulty] == hashString {
		return true
	}
	return false
}

// newBlock add's a new block to the chain and resets the transactions as new transactions will be added
// to the next block
func (bc *Blockchain) newBlock(proof int64, previousHash string) Block {
	prevHash := previousHash
	if previousHash == "" {
		prevBlock := bc.Chain[len(bc.Chain)-1]
		prevHash = hash(prevBlock)
	}

	block := Block{
		Index:        int64(len(bc.Chain) + 1),
		Timestamp:    time.Now().UnixNano(),
		Transactions: bc.Transactions,
		Proof:        proof,
		PreviousHash: prevHash,
	}

	bc.Transactions = nil // reset transactions as the block will be added to the chain
	bc.Chain = append(bc.Chain, block)
	return block
}

// initBlockchain initialises the blockchain
// Returns a pointer to the blockchain object that the app can alter later on
func initBlockchain() *Blockchain {
	// init the blockchain
	newBlockchain := &Blockchain{
		Chain:        make([]Block, 0),
		Transactions: make([]Transaction, 0),
		Nodes:        nil,
	}
	fmt.Printf("init Blockchain\n %v\n", newBlockchain)

	// adding a first, Genesis, Block to the Chain
	b := newBlockchain.newBlock(100, "_")
	fmt.Printf("adding a Block:\n %v\n", b)
	fmt.Printf("Blockchain:\n %v\n", newBlockchain)
	return newBlockchain // pointer
}

// validate. Determines if a given blockchain is valid.
// Returns bool, true if valid
func (bc *Blockchain) validate() bool {

	chainLength := len(bc.Chain)
	fmt.Printf("Validating a chain with a chainLength of %d\n", chainLength)

	if chainLength == 1 {
		fmt.Println("chain has only one block yet, thus  valid")
		return true
	}

	for i := 1; i < chainLength; i++ {
		//# Check that the hash of the block is correct
		//if block['previous_hash'] != self.hash(last_block):
		//return False

		previous := bc.Chain[i - 1]
		current := bc.Chain[i]

		if current.PreviousHash != hash(previous) {
			fmt.Println("invalid hash")
			fmt.Printf("Previous block: %d\n", previous.Index)
			fmt.Printf("Current block: %d\n", current.Index)
			return false
		}

		//# Check that the Proof of Work is correct
		//if not self.valid_proof(last_block['proof'], block['proof']):
		//return False

		if !bc.validProof(previous.Proof, current.Proof) {
			fmt.Println("invalid proof")
			fmt.Printf("Previous block: %d\n", previous.Index)
			fmt.Printf("Current block: %d\n", current.Index)
			return false
		}
	}


	return true

	//last_block = chain[0]
	//current_index = 1
	//
	//while current_index < len(chain):
	//block = chain[current_index]
	//print(f'{last_block}')
	//print(f'{block}')
	//print("\n-----------\n")
	//# Check that the hash of the block is correct
	//if block['previous_hash'] != self.hash(last_block):
	//return False
	//
	//# Check that the Proof of Work is correct
	//if not self.valid_proof(last_block['proof'], block['proof']):
	//return False
	//
	//last_block = block
	//current_index += 1
	//
	//return True
}
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/grrrben/golog"
	"net/http"
	"time"
)

// how many 0's do we want to check
const hashDifficulty int8 = 4

type Blockchain struct {
	Chain        []Block
	Transactions []Transaction
}

type chainService interface {
	newBlock() bool
	newTransaction() bool
	hash(Block) string
	lastBlock() Block
	proofOfWork(lastProof int64) int64
	validProof(proof int64, lastProof int64) bool
	validate() bool
	resolve() bool
}

// newTransaction will create a Transaction to go into the next Block to be mined.
// The Transaction is stored in the Blockchain obj.
// Returns (int) the Index of the Block that will hold this Transaction
func (bc *Blockchain) newTransaction(tr Transaction) int64 {
	bc.Transactions = append(bc.Transactions, tr)
	return bc.lastBlock().Index + 1
}

// Hash Creates a SHA-256 hash of a Block
func hash(b Block) string {
	if debug {
		golog.Infof("hashing block %d\n", b.Index)
	}

	// Data for binary.Write must be a fixed-size value or a slice of fixed-size values,
	// or a pointer to such data.
	// @todo Marshalling the struct to json is a workaround... But it works
	// @todo might be able to fix it with a char(length) instead of string?
	jsonblock, errr := json.Marshal(b)
	if errr != nil {
		if debug {
			golog.Errorf("Error: %s", errr)
		}
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, jsonblock)
	if err != nil {
		if debug {
			golog.Errorf("Could not compute hash: %s", err)
		}
	}
	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())) // %x; base 16, with lower-case letters for a-f
}

// lastBlock returns the last Block in the Chain
func (bc *Blockchain) lastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) proofOfWork(lastProof int64) int64 {
	// Simple Proof of Work Algorithm:
	// - Find a number p' such that hash(lp') contains leading 4 zeroes, where
	// - l is the previous Proof, and p' is the new Proof
	var proof int64 = 0
	i := 0
	for !bc.validProof(lastProof, proof) {
		proof += 1
		i++
	}
	if debug {
		golog.Infof("Proof found in %d cycles (difficulty %d)\n", i, hashDifficulty)
	}
	return proof

}

// validProof is called until it finds an acceptable hash and returns true
func (bc *Blockchain) validProof(proof int64, lastProof int64) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := fmt.Sprintf("%x", sha256.Sum256([]byte(guess)))

	var i int8
	hashString := ""
	for i = 0; i < hashDifficulty; i++ {
		// todo move this out of the loopt
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
	}
	if debug {
		golog.Infof("init Blockchain\n %v\n", newBlockchain)
	}
	// adding a first, Genesis, Block to the Chain
	b := newBlockchain.newBlock(100, "_")
	if debug {
		golog.Infof("adding a Block:\n %v\n", b)
	}
	return newBlockchain // pointer
}

// validate. Determines if a given blockchain is valid.
// Returns bool, true if valid
func (bc *Blockchain) validate() bool {

	chainLength := len(bc.Chain)
	if debug {
		golog.Infof("Validating a chain with a chainLength of %d\n", chainLength)
	}

	if chainLength == 1 {
		if debug {
			golog.Info("chain has only one block yet, thus  valid")
		}
		return true
	}

	for i := 1; i < chainLength; i++ {
		// Check that the hash of the block is correct
		// if block['previous_hash'] != self.Hash(last_block):
		// return False
		previous := bc.Chain[i-1]
		current := bc.Chain[i]

		if current.PreviousHash != hash(previous) {
			if debug {
				golog.Info("invalid Hash")
				golog.Infof("Previous block: %d\n", previous.Index)
				golog.Infof("Current block: %d\n", current.Index)
			}
			return false
		}

		// Check that the Proof of Work is correct
		// if not self.valid_proof(last_block['proof'], block['proof']):
		// return False
		if !bc.validProof(previous.Proof, current.Proof) {
			if debug {
				golog.Info("invalid proof")
				golog.Infof("Previous block: %d\n", previous.Index)
				golog.Infof("Current block: %d\n", current.Index)
			}
			return false
		}
	}
	return true
}

// resolve is the Consensus Algorithm, it resolves conflicts
// by replacing our chain with the longest one in the network.
// Returns bool. True if our chain was replaced, false if not
func (bc *Blockchain) resolve() bool {
	golog.Infof("Resolving conflicts (clients %d):", len(cls.List))
	length := len(bc.Chain)
	replaced := false
	for _, cl := range cls.List {
		if cl == me {
			continue
		}
		url := fmt.Sprintf("%s%s:%d/chain", cl.Protocol, cl.Ip, cl.Port)
		golog.Infof("%s\n", url)

		resp, err := http.Get(url)
		if err != nil {
			golog.Warningf("Chain request error: %s", err)
			// I don't want to panic here, but it could be a good idea to
			// remove the client from the list
			continue
		}

		var extChain Blockchain
		decodingErr := json.NewDecoder(resp.Body).Decode(&extChain)
		defer resp.Body.Close()

		if decodingErr != nil {
			golog.Warningf("Could not decode JSON of external blockchain\n")
			golog.Warningf("Error: %s\n", err)
			continue
		}

		if len(extChain.Chain) > length {
			golog.Infof("Found a new blockchain with length %d.\n", len(extChain.Chain))
			golog.Infof("Our blockchain had a length of %d.\n", length)
			golog.Infof("Blockchain replaced.")

			// it might be better to fetch a list of all client's chain length first, then replace ours
			// with the largest one.
			bc.Chain = extChain.Chain
			replaced = true
		}
	}
	return replaced
}

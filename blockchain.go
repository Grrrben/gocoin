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

// how many zero's do we want in the hash
const hashDifficulty int8 = 4

// This should be the hash ending in the proof of work
const hashEndsWith string = "0000"

type Blockchain struct {
	Chain        []Block
	Transactions []Transaction
}

// newTransaction will create a Transaction to go into the next Block to be mined.
// The Transaction is stored in the Blockchain obj.
// Returns (int) the Index of the Block that will hold this Transaction
func (bc *Blockchain) newTransaction(tr Transaction) int64 {
	tr.Time = time.Now().UnixNano()
	bc.Transactions = append(bc.Transactions, tr)
	return bc.lastBlock().Index + 1
}

// Hash Creates a SHA-256 hash of a Block
func hash(bl Block) string {
	golog.Infof("hashing block %d\n", bl.Index)

	// Data for binary.Write must be a fixed-size value or a slice of fixed-size values,
	// or a pointer to such data.
	// @todo Marshalling the struct to json is a workaround... But it works
	// @todo might be able to fix it with a char(length) instead of string?
	jsonblock, errr := json.Marshal(bl)
	if errr != nil {
		golog.Errorf("Error: %s", errr)
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, jsonblock)
	if err != nil {
		golog.Errorf("Could not compute hash: %s", err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())) // %x; base 16, with lower-case letters for a-f
}

// lastBlock returns the last Block in the Chain
func (bc *Blockchain) lastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) proofOfWork(lastProof int64) int64 {
	// Simple Proof of Work Algorithm:
	// - Find a number p' such that hash(lp') contains leading X zeroes, where
	// - l is the previous Proof, and p' is the new Proof
	var proof int64 = 0
	i := 0
	for !bc.validProof(lastProof, proof) {
		proof += 1
		i++
	}
	golog.Infof("Proof found in %d cycles (difficulty %s)\n", i, hashEndsWith)
	return proof

}

// validProof is called until it finds an acceptable hash and returns true
func (bc *Blockchain) validProof(proof int64, lastProof int64) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := fmt.Sprintf("%x", sha256.Sum256([]byte(guess)))

	if guessHash[:hashDifficulty] == hashEndsWith {
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
	cls.announceMinedBlocks(block)
	return block
}

// addBlock performs a validity check on the new block, if valid it add's the block to the chain.
// Return bool
func (bc *Blockchain) addBlock(bl Block) bool {

	lastBlock := bc.Chain[len(bc.Chain)-1]

	if bc.validProof(lastBlock.Proof, bl.Proof) {
		golog.Info("Added a new block due to an announcement.")
		bc.Chain = append(bc.Chain, bl)
		return true
	}
	golog.Warning("Could not add the newly announced block.")
	return false
}

// analyseInvalidBlock
// shows us why a newly sent block could not be added to the chain.
// and tries to add more blocks if we are missing multiple.
func (bc *Blockchain) analyseInvalidBlock(bl Block, sender string) bool {

	lastBlock := bc.Chain[len(bc.Chain)-1]

	golog.Info("----------------------------------")
	defer golog.Info("----------------------------------")
	golog.Infof("Analysing block: index: %d", bl.Index)
	golog.Infof("%v", bl)
	golog.Infof("Last block: index: %d", lastBlock.Index)
	golog.Infof("%v", lastBlock)

	if lastBlock.Index < (bl.Index - 1) {
		var i int64 // 0
		for {
			i++
			var nextBlock Block

			url := fmt.Sprintf("%s/block/index/%d", sender, lastBlock.Index+i)
			golog.Infof("Fetching block %d from $s", lastBlock.Index+i, sender)

			resp, err := http.Get(url)
			if err != nil {
				golog.Warningf("Request error: %s", err)
				return false
			}

			decodingErr := json.NewDecoder(resp.Body).Decode(&nextBlock)
			if decodingErr != nil {
				golog.Warningf("Decoding error: %s", err)
				return false
			}

			success := bc.addBlock(nextBlock)
			if success == false {
				golog.Warningf("Could not add block %d from %s", lastBlock.Index+i, sender)
				return false
			}
			defer resp.Body.Close()

			if (lastBlock.Index + i) == bl.Index {
				golog.Infof("Successfully added %d blocks", i)
				break
			}
		}
	} else {
		// something else went wrong.
		golog.Warning("Unable to analyse")
		return false
	}

	return true
}

// initBlockchain initialises the blockchain
// Returns a pointer to the blockchain object that the app can alter later on
func initBlockchain() *Blockchain {
	// init the blockchain
	newBlockchain := &Blockchain{
		Chain:        make([]Block, 0),
		Transactions: make([]Transaction, 0),
	}
	golog.Infof("init Blockchain\n %v", newBlockchain)

	if me.Port == 8000 {
		// Mother node. Adding a first, Genesis, Block to the Chain
		b := newBlockchain.newBlock(100, zerohash)
		golog.Infof("Adding Genesis Block:\n %v", b)
	} else {
		newBlockchain.resolve()
		golog.Infof("Resolving the blockchain")
	}

	return newBlockchain // pointer
}

// validate. Determines if a given blockchain is valid.
// Returns bool, true if valid
func (bc *Blockchain) validate() bool {

	chainLength := len(bc.Chain)

	if chainLength == 1 {
		return true
	}

	for i := 1; i < chainLength; i++ {
		// Check that the hash of the block is correct
		// if block['previous_hash'] != self.Hash(last_block):
		// return False
		previous := bc.Chain[i-1]
		current := bc.Chain[i]

		if current.PreviousHash != hash(previous) {
			golog.Warning("invalid Hash")
			golog.Warningf("Previous block: %d\n", previous.Index)
			golog.Warningf("Current block: %d\n", current.Index)
			return false
		}

		// Check that the Proof of Work is correct
		// if not self.valid_proof(last_block['proof'], block['proof']):
		// return False
		if !bc.validProof(previous.Proof, current.Proof) {
			golog.Warning("invalid proof")
			golog.Warningf("Previous block: %d\n", previous.Index)
			golog.Warningf("Current block: %d\n", current.Index)
			return false
		}
	}
	return true
}

// mine Mines a block and puts all transactions in the block
// An incentive is paid to the miner and the list of transactions is cleared
func (bc *Blockchain) mine() Block {
	lastBlock := bc.lastBlock()
	lastProof := lastBlock.Proof

	proof := bc.proofOfWork(lastProof)
	tr := Transaction{zerohash, me.Hash, 1, time.Now().UnixNano()}
	bc.newTransaction(tr)
	block := bc.newBlock(proof, "")
	return block
}

// resolve is the Consensus Algorithm, it resolves conflicts
// by replacing our chain with the longest one in the network.
// Returns bool. True if our chain was replaced, false if not
func (bc *Blockchain) resolve() bool {
	golog.Infof("Resolving conflicts (clients %d):", len(cls.List))
	replaced := false
	for _, cl := range cls.List {
		if cl == me {
			continue
		}
		url := fmt.Sprintf("%s/chain", cls.getAddress(cl))
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
			golog.Warningf("Could not decode JSON of external blockchain: %s", err)
			continue
		}

		if len(extChain.Chain) > len(bc.Chain) {
			golog.Infof("Blockchain replaced. Found length of %d instead of current %d.", len(extChain.Chain), len(bc.Chain))
			// it might be better to fetch a list of all client's chain length first, then replace ours
			// with the largest one.
			bc.Chain = extChain.Chain
			replaced = true
		}
	}
	return replaced
}

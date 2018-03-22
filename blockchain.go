package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/grrrben/glog"
)

// how many zero's do we want in the hash
const hashDifficulty int8 = 4

// This should be the hash ending in the proof of work
const hashEndsWith string = "0000"

// The incentive paid to the miner of a minted block
const minersIncentive = 1

type Blockchain struct {
	Chain        []Block
	Transactions []Transaction
}

// StatusReport is used to fetch the information regarding the blockchain from other nodes in the network.
type StatusReport struct {
	Length int `json:"length"`
}

// NodeLength represents the length of the blockchain of a particular node.
type NodeLength struct {
	node   Node
	length int
}

// newTransaction will create a Transaction to go into the next Block to be mined.
// The Transaction is stored in the Blockchain obj.
// Returns the Transation with an added Time property
func (bc *Blockchain) newTransaction(transaction Transaction) (tr Transaction, err error) {
	_, err = checkTransaction(transaction)

	if err != nil {
		return transaction, err
	} else {
		if transaction.Time == 0 {
			transaction.Time = time.Now().UnixNano()
		}
		bc.Transactions = append(bc.Transactions, transaction)
		return transaction, nil
	}
}

// isNonExistingTransaction loops the current list of Transactions
// to check if the new Transactions is already known on this Node
func (bc *Blockchain) isNonExistingTransaction(newTr Transaction) bool {
	for _, existingTr := range bc.Transactions {
		if checkHashesEqual(newTr, existingTr) {
			return false
		}
	}
	return true
}

// clearTransactions loops all transactions in this node and filters out all transactions that are
// persisted in the mined block
func (bc *Blockchain) clearTransactions(trs []Transaction) {
	// get a map of all hashes and their corresponding Transactions
	var hashesInBlock = map[string]Transaction{}
	for _, tr := range trs {
		hashesInBlock[tr.getHash()] = tr
	}

	// Set the transactions not found in the announced block to this chain's transaction List
	var transactionsNotInMinedBlock []Transaction
	for _, tr := range bc.Transactions {
		_, exists := hashesInBlock[tr.getHash()]
		if !exists {
			glog.Infof("Transaction does not exist, keeping it:\n %v", tr)
			transactionsNotInMinedBlock = append(transactionsNotInMinedBlock, tr)
		}
	}
	bc.Transactions = transactionsNotInMinedBlock
}

// Hash Creates a SHA-256 hash of a Block
func hash(bl Block) string {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(bl)
	if err != nil {
		glog.Errorf("Could not compute hash: %s", err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())) // %x; base 16, with lower-case letters for a-f
}

// lastBlock returns the last Block in the Chain
func (bc *Blockchain) lastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

// proofOfWork is a simple Proof of Work Algorithm:
// Find a number p such that hash('lp') contains leading X zeroes, where
// l is the previous Proof, and p is the new Proof
func (bc *Blockchain) proofOfWork(lastProof int64) int64 {
	var proof int64 = 0
	for !bc.validProof(lastProof, proof) {
		proof++
	}
	glog.Infof("Proof found in %d cycles (difficulty %d)\n", proof, hashDifficulty)
	return proof
}

// validProof is called until it finds an acceptable hash and returns true
func (bc *Blockchain) validProof(proof, lastProof int64) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := fmt.Sprintf("%x", sha256.Sum256([]byte(guess)))
	return guessHash[:hashDifficulty] == hashEndsWith
}

// newBlock add's a new block to the chain and resets the transactions as new transactions will be added
// to the next block
func (bc *Blockchain) newBlock(proof int64) Block {
	var prevHash string
	if len(bc.Chain) == 0 {
		// this is the genesis block
		prevHash = zerohash
	} else {
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
	nodes.announceMinedBlocks(block)
	return block
}

// addBlock performs a validity check on the new block, if valid it add's the block to the chain.
// Return bool
func (bc *Blockchain) addBlock(bl Block) (Block, error) {

	lastBlock := bc.Chain[len(bc.Chain)-1]

	if bc.validProof(lastBlock.Proof, bl.Proof) {
		glog.Info("Added a new block due to an announcement.")
		bc.Chain = append(bc.Chain, bl)
		return bl, nil
	}
	return bl, errors.New("Could not add the newly announced block.")
}

// analyseInvalidBlock
// shows us why a newly sent block could not be added to the chain.
// and tries to add more blocks if we are missing multiple.
func (bc *Blockchain) analyseInvalidBlock(bl Block, sender string) bool {

	lastBlock := bc.Chain[len(bc.Chain)-1]

	glog.Info("----------------------------------")
	defer glog.Info("----------------------------------")
	glog.Infof("Analysing block: index: %d", bl.Index)
	glog.Infof("%v", bl)
	glog.Infof("Last block: index: %d", lastBlock.Index)
	glog.Infof("%v", lastBlock)

	if lastBlock.Index < (bl.Index - 1) {
		var i int64 // 0
		for {
			i++
			var nextBlock Block

			url := fmt.Sprintf("%s/block/index/%d", sender, lastBlock.Index+i)
			glog.Infof("Fetching block %d from $s", lastBlock.Index+i, sender)

			resp, err := http.Get(url)
			if err != nil {
				glog.Warningf("Request error: %s", err)
				return false
			}

			decodingErr := json.NewDecoder(resp.Body).Decode(&nextBlock)
			if decodingErr != nil {
				glog.Warningf("Decoding error: %s", err)
				return false
			}

			_, err = bc.addBlock(nextBlock)
			if err != nil {
				glog.Warningf("Could not add block %d from %s: %s", lastBlock.Index+i, sender, err.Error())
				return false
			}
			defer resp.Body.Close()

			if (lastBlock.Index + i) == bl.Index {
				glog.Infof("Successfully added %d blocks", i)
				break
			}
		}
	} else {
		// something else went wrong.
		glog.Warning("Unable to analyse")
		return false
	}

	return true
}

// initBlockchain initialises the blockchain
// Returns a pointer to the blockchain object that the app can alter later on
// If there already is a network, the chain is fetched from the network, otherwise a genesis block is created.
func initBlockchain() *Blockchain {
	// init the blockchain
	newBlockchain := &Blockchain{
		Chain:        make([]Block, 0),
		Transactions: make([]Transaction, 0),
	}
	glog.Infof("init Blockchain\n %v", newBlockchain)

	if me.Port == 8000 {
		// Mother node. Adding a first, Genesis, Block to the Chain
		b := newBlockchain.newBlock(100)
		glog.Infof("Adding Genesis Block:\n %v", b)
	} else {
		newBlockchain.resolve()
		glog.Infof("Resolving the blockchain")
	}

	return newBlockchain // pointer
}

// getCurrentTransactions get's the transactions from other nodes.
// it is used at the startup
func (bc *Blockchain) getCurrentTransactions() bool {
	defer glog.Flush()
	if len(nodes.List) > 1 {
		for _, node := range nodes.List {
			url := fmt.Sprintf("%s/transactions", node.getAddress())

			if me.getAddress() == node.getAddress() {
				// it is I, skip it
				continue
			}
			resp, err := http.Get(url)
			if err != nil {
				glog.Warningf("Transactions request error: %s", err)
				continue
			}

			var transactions []Transaction

			decodingErr := json.NewDecoder(resp.Body).Decode(&transactions)

			if decodingErr != nil {
				glog.Warningf("Could not decode JSON of external transactions: %s", err)
				continue
			}
			resp.Body.Close()
			glog.Infof("Found %d transactions on another node.", len(transactions))
			bc.Transactions = transactions
			return true
		}
		glog.Warning("No transactions found on other nodes")
	}
	glog.Info("First node. No transactions added")
	return false
}

// validate. Determines if a given blockchain is valid.
func (bc *Blockchain) validate() bool {
	defer glog.Flush()
	chainLength := len(bc.Chain)

	if chainLength == 1 {
		return true
	}

	for i := 1; i < chainLength; i++ {
		previous := bc.Chain[i-1]
		current := bc.Chain[i]

		// Check that the hash of the block is correct
		if current.PreviousHash != hash(previous) {
			glog.Warningf("Invalid Hash in blockchain, block %d cannot be placed before block %d", previous.Index, current.Index)
			return false
		}

		// Check that the Proof of Work is correct
		if !bc.validProof(previous.Proof, current.Proof) {
			glog.Warningf("Invalid proof of block %d, with previous block %d", current.Index, previous.Index)
			return false
		}
	}
	return true
}

// mine Mines a block and puts all transactions in the block
// An incentive is paid to the miner and the list of transactions is cleared
func (bc *Blockchain) mine() (Block, error) {
	var block Block
	lastBlock := bc.lastBlock()
	lastProof := lastBlock.Proof

	proof := bc.proofOfWork(lastProof)
	transaction := Transaction{
		zerohash,
		me.Hash,
		minersIncentive,
		fmt.Sprintf("Mined by %s", me.getAddress()),
		time.Now().UnixNano(),
	}
	_, err := bc.newTransaction(transaction)
	if err != nil {
		return block, err
	}
	block = bc.newBlock(proof)
	return block, nil
}

// resolve is the Consensus Algorithm, it resolves conflicts by replacing our chain with the longest one in the network.
// Returns bool. True if our chain was replaced, false if not
func (bc *Blockchain) resolve() bool {
	glog.Infof("Resolving conflicts (nodes %d):", len(nodes.List))
	replaced := false

	// first, let's grep some of the lengths of the different node chains.
	nodes := bc.chainLengthPerNode()

	for _, pair := range nodes {

		var node Node
		node = pair.Key.(Node) // getting the type back as the interface{} signature didn't give hints
		if node == me {
			continue
		}
		url := fmt.Sprintf("%s/chain", node.getAddress())
		resp, err := http.Get(url)
		if err != nil {
			glog.Warningf("Chain request error: %s", err)
			// I don't want to panic here, but it could be a good idea to
			// remove the node from the list
			continue
		}

		var extChain Blockchain
		decodingErr := json.NewDecoder(resp.Body).Decode(&extChain)

		if decodingErr != nil {
			glog.Warningf("Could not decode JSON of external blockchain: %s", err)
			continue
		}

		if len(extChain.Chain) > len(bc.Chain) {
			// check if the chain is valid.
			oldChain := bc.Chain
			bc.Chain = extChain.Chain
			valid := bc.validate()

			if valid {
				glog.Infof("Blockchain replaced. Found length of %d instead of current %d.", len(extChain.Chain), len(bc.Chain))
				glog.Infof("Synced with %s\n", node.getAddress())
				replaced = true
			} else {
				// reset to old blockchain
				bc.Chain = oldChain
			}

		}
		resp.Body.Close()

		if replaced {
			// we have a new, valid, chain.
			break
		}

	}
	return replaced
}

// chainLengthPerNode get a map of nodes with their respective chain length
func (bc *Blockchain) chainLengthPerNode() PairList {
	// a map of Nodes with their chain length, the interface{} is used as a key so it is compatible with the sortMapDescending function
	nodeLength := make(map[interface{}]int)
	// a channel with the cl vs length struct
	nodeChannel := make(chan NodeLength, 10)
	// in case something goes wrong, show a couple of errors
	errChannel := make(chan error, 4)

	var wg sync.WaitGroup

	for i, cl := range nodes.List {
		if cl == me {
			continue
		}
		wg.Add(1)
		go chainLengthOfNode(cl, &wg, nodeChannel, errChannel)
		if i > 10 {
			break // max 10, but sooner if less nodes are connected
		}
	}

	// wait for the sync.WaitGroup to be completed, afterwards the channels can be closed safely
	wg.Wait()
	close(nodeChannel)
	close(errChannel)

	for receiver := range nodeChannel {
		glog.Infof("received: %v", receiver)
		nodeLength[receiver.node] = receiver.length
	}

	for err := range errChannel {
		glog.Warningf("Error in fetching list of node statusses", err.Error())
	}

	glog.Infof("Length of nodes:\n%v\n", len(nodeLength))
	// watch it; When iterating over a map with a range loop, the iteration order is not specified and is not
	// guaranteed to be the same from one iteration to the next. Thus, sort it first.
	return sortMapDescending(nodeLength)
}

// chainLengthOfNode Goroutine. Helper function that collects information from nodes and puts it in the channel
func chainLengthOfNode(cl Node, wg *sync.WaitGroup, channel chan NodeLength, errorChannel chan error) {
	var report StatusReport
	defer wg.Done()

	url := fmt.Sprintf("%s/status", cl.getAddress())
	resp, err := http.Get(url)
	if err != nil {
		select {
		case errorChannel <- err:
			// first X errors to this channel
		default:
			// ok
		}
	} else {
		// We have no error, thus we can decode the response in the repost val.
		decodingErr := json.NewDecoder(resp.Body).Decode(&report)

		if decodingErr != nil {
			select {
			case errorChannel <- decodingErr:
				// first X errors to this channel
			default:
				// ok
			}
		} else {
			defer resp.Body.Close()
		}

		clen := NodeLength{cl, report.Length}
		channel <- clen
	}
}

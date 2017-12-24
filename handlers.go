package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grrrben/golog"
	"net/http"
	"strconv"
)

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Hello world")
}

// wallet Shows some stats of a wallet, including the credits available
func (a *App) wallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	resp := map[string]interface{}{
		"success": true,
		"credit":  getWalletCredits(hash),
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// transactions shows all transactions made by a wallet with {hash}
// {hash} is given by POST data from the call
func (a *App) transactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	transactions := []Transaction{}

	// check all blocks, see if the hash is the sender or receiver.
	for _, block := range bc.Chain {
		for _, transaction := range block.Transactions {
			if transaction.Sender == hash || transaction.Recipient == hash {
				transactions = append(transactions, transaction)
			}
		}
	}
	respondWithJSON(w, http.StatusOK, transactions)
}

// currentTransactions shows all transactions that are not in a block yet
func (a *App) currentTransactions(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, bc.Transactions)
}

// distributedTransaction receives a transaction from another client in the network.
// It is used to distribute the _unmined_ transactions throughout the network
func (a *App) distributedTransaction(w http.ResponseWriter, r *http.Request) {
	defer golog.Flush()
	golog.Infof("starting distributedTransaction on Client: %s", me.getAddress())

	type Payload struct {
		Transaction Transaction `json:"transaction"`
		Sender      string      `json:"sender"`
	}

	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		golog.Warningf("Invalid Transaction (Unable to decode) on Client: %s", me.getAddress())
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid Transaction (Unable to decode)")
	} else {
		golog.Infof("payload: %v", payload)
		if bc.isNonExistingTransaction(payload.Transaction) {
			golog.Infof("transaction: %v", payload.Transaction)
			_, err = bc.newTransaction(payload.Transaction)
			if err != nil {
				golog.Warningf("%s on Client: %s", err.Error(), me.getAddress())
				respondWithError(w, http.StatusUnprocessableEntity, err.Error())
			} else {
				golog.Infof("Transaction added on Client: %s", me.getAddress())
				respondWithJSON(w, http.StatusOK, "Transaction added")
			}
		} else {
			golog.Warningf("Invalid Transaction (Already exists) on Client: %s", me.getAddress())
			respondWithError(w, http.StatusUnprocessableEntity, "Invalid Transaction (Already exists)")
		}

	}
}

// newTransaction adds a transaction, which consists of:
// Sender string
// Recipient string
// Amount float32
func (a *App) newTransaction(w http.ResponseWriter, r *http.Request) {
	var tr Transaction
	err := json.NewDecoder(r.Body).Decode(&tr)

	if err != nil {
		golog.Warningf("Invalid Transaction (%s)", err.Error())
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid Transaction (Unable to decode)")
	} else {
		// for a new transaction, time should be added by the system
		tr.Time = 0
		addedTransaction, err := bc.newTransaction(tr)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		} else {
			// all OK. Add the transaction and distribute it.
			cls.distributeTransaction(addedTransaction) // distribution
			respondWithJSON(w, http.StatusOK, "Transaction added")
		}
	}
}

// distributedBlock is a receiver for blocks mined by other clients.
// It catches the newly mined block and checks for validity on his own chain
// If it is valid the block is added and a statusOk is returned.
// Otherwise it gives an error
func (a *App) distributedBlock(w http.ResponseWriter, r *http.Request) {
	// fetching the block that came with the request
	decoder := json.NewDecoder(r.Body)

	type Payload struct {
		NewBlock Block  `json:"block"`
		Sender   string `json:"sender"`
	}

	var payload Payload
	err := decoder.Decode(&payload)
	if err != nil {
		golog.Warning("Could not decode postdata of new block")
		respondWithError(w, http.StatusBadRequest, "invalid json")
		panic(err)
	}
	block, err := bc.addBlock(payload.NewBlock)

	if err == nil {
		bc.clearTransactions(block.Transactions)
		resp := map[string]interface{}{
			"success": true,
			"message": "New block added",
		}
		respondWithJSON(w, http.StatusOK, resp)
	} else {
		repair := bc.analyseInvalidBlock(payload.NewBlock, payload.Sender)

		if repair == false {
			// better resolve..?
			respondWithError(w, http.StatusConflict, "Invalid block")
		} else {
			resp := map[string]interface{}{
				"success": true,
				"message": "New blocks added",
			}
			respondWithJSON(w, http.StatusOK, resp)
		}
	}
}

// resolve Resolving conflict between chains in the network
func (a *App) resolve(w http.ResponseWriter, r *http.Request) {
	resolved := bc.resolve()
	respondWithJSON(w, http.StatusOK, resolved)
}

// lastblock Serves single block
func (a *App) lastblock(w http.ResponseWriter, r *http.Request) {
	block := bc.Chain[len(bc.Chain)-1]
	resp := map[string]interface{}{"success": true, "block": block}
	respondWithJSON(w, http.StatusOK, resp)
}

// block Serves single block identified by it's hash
func (a *App) block(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	found := false

	for _, bl := range bc.Chain {
		if bl.PreviousHash == hash {
			found = true
			resp := map[string]interface{}{"success": true, "block": bl}
			respondWithJSON(w, http.StatusOK, resp)
		}
		break
	}

	if found == false {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not find block by hash %s", hash))
	}
}

// blockByIndex Serves single block identified by it's index
func (a *App) blockByIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawIndex := vars["index"]

	index, err := strconv.ParseInt(rawIndex, 10, 16) // always gives an uint64...
	if err != nil {
		golog.Errorf("Unable to cast block Index %s to int: %s", rawIndex, err)
	}

	found := false

	for _, bl := range bc.Chain {
		if bl.Index == index {
			found = true
			resp := map[string]interface{}{"success": true, "block": bl}
			respondWithJSON(w, http.StatusOK, resp)
		}
		break
	}

	if found == false {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not find block by index %s", index))
	}
}

// chainStatus
func (a *App) chainStatus(w http.ResponseWriter, r *http.Request) {
	hash := bc.Chain[len(bc.Chain)-1].PreviousHash
	resp := map[string]interface{}{"length": len(bc.Chain), "hash": hash}
	respondWithJSON(w, http.StatusOK, resp)
}

// connectClient Connect a Client to the network which is represented
// in the Clients.list The postdata should consist of a standard Client
func (a *App) connectClient(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newCl Client
	err := decoder.Decode(&newCl)
	if err != nil {
		golog.Warning("Could not decode postdata of new client")
		respondWithError(w, http.StatusBadRequest, "invalid json")
		panic(err)
	}
	// register the client
	added := cls.addClient(newCl)
	if added {
		resp := map[string]interface{}{"Client": newCl, "total": cls.num()}
		respondWithJSON(w, http.StatusOK, resp)
	} else {
		respondWithError(w, http.StatusConflict, "Client could not be added")
	}
}

// getClients response is the list of Clients
func (a *App) getClients(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{"list": cls.List, "length": len(cls.List)}
	respondWithJSON(w, http.StatusOK, resp)
}

// chain shows the entire blockchain
func (a *App) chain(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{"chain": bc.Chain, "transactions": bc.Transactions, "length": len(bc.Chain)}
	respondWithJSON(w, http.StatusOK, resp)
}

// validate checks the entire blockchain
func (a *App) validate(w http.ResponseWriter, r *http.Request) {
	isValid := bc.validate()
	resp := map[string]interface{}{"valid": isValid, "length": len(bc.Chain)}
	respondWithJSON(w, http.StatusOK, resp)
}

// mine Mines a block and puts all transactions in the block
// An incentive is paid to the miner and the list of transactions is cleared
func (a *App) mine(w http.ResponseWriter, r *http.Request) {
	block, err := bc.mine()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		resp := map[string]interface{}{
			"message":      "New block mined.",
			"Block":        block,
			"length":       len(bc.Chain),
			"transactions": len(block.Transactions),
		}
		respondWithJSON(w, http.StatusOK, resp)
	}
}

// GetConfig test of the config needs to be loaded and returns the Config file.
func GetConfig() Config {
	var cf *Config
	if cf == nil {
		readConfig()
	}
	return config
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

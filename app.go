package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"os"
)

type Server struct {
	Server string
	Port   int
}

type Config struct {
	Server Server
}

var config Config

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() {
	config = GetConfig()
	bc = initBlockchain()

	if debug {
		fmt.Println("Starting with a base blockchain:")
		fmt.Printf("Blockchain:\n %v\n", bc)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run() {
	port := fmt.Sprintf("%d", config.Server.Port)
	fmt.Println("Starting server")
	fmt.Printf("Running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.index).Methods("GET")
	a.Router.HandleFunc("/transaction", a.newTransaction).Methods("POST")
	a.Router.HandleFunc("/mine", a.mine).Methods("GET")
	a.Router.HandleFunc("/chain", a.chain).Methods("GET")
	a.Router.HandleFunc("/validate", a.validate).Methods("GET")
}

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Hello world")
}

func (a *App) newTransaction(w http.ResponseWriter, r *http.Request) {
	// Sender string
	// Recipient string
	// Amount float64
	var tr Transaction

	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Transaction")
	}

	bc.newTransaction(tr)

	respondWithJSON(w, http.StatusOK, "Transaction added")
}

func (a *App) chain(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{"chain": bc.Chain, "transactions": bc.Transactions, "length": len(bc.Chain)}
	respondWithJSON(w, http.StatusOK, resp)
}

func (a *App) validate(w http.ResponseWriter, r *http.Request) {
	isValid := bc.validate()
	resp := map[string]interface{}{"valid": isValid, "length": len(bc.Chain)}
	respondWithJSON(w, http.StatusOK, resp)
}

func (a *App) mine(w http.ResponseWriter, r *http.Request) {
	//# We run the proof of work algorithm to get the next proof...
	//last_block = blockchain.last_block
	//last_proof = last_block['proof']
	//proof = blockchain.proof_of_work(last_proof)
	//
	//# We must receive a reward for finding the proof.
	//# The sender is "0" to signify that this node has mined a new coin.
	//blockchain.new_transaction(
	//sender="0",
	//recipient=node_identifier,
	//amount=1,
	//)
	//
	//# Forge the new Block by adding it to the chain
	//block = blockchain.new_block(proof)
	//
	//response = {
	//'message': "New Block Forged",
	//'index': block['index'],
	//'transactions': block['transactions'],
	//'proof': block['proof'],
	//'previous_hash': block['previous_hash'],
	//}
	//return jsonify(response), 200

	lastBlock := bc.lastBlock()
	lastProof := lastBlock.Proof

	proof := bc.proofOfWork(lastProof)
	tr := Transaction{"0", "recipient", 1}
	bc.newTransaction(tr)
	block := bc.newBlock(proof, "")

	resp := map[string]interface{}{
		"message":      "New block mined.",
		"Block":        block,
		"length":       len(bc.Chain),
		"transactions": len(block.Transactions),
	}
	respondWithJSON(w, http.StatusOK, resp)
}

// GetConfig test of the config needs to be loaded and returns the Config file.
func GetConfig() Config {
	var cf *Config
	if cf == nil {
		readConfig()
	}
	return config
}

// readConfig loads the config file.
// It tests for the existence of the file and whether or not it can be decoded by a TOML decoder
func readConfig() {
	var configfile = "config.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
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

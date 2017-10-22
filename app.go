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

func (a *App) Initialize(port uint16) {
	config = GetConfig()
	bc = initBlockchain()

	messenger("Starting with a base blockchain:")
	messenger("Blockchain:\n %v\n", bc)

	// add the Client to the stack
	cls = initClients() // a pointer to the Clients struct
	cl := Client{
		Protocol: "http://",
		Ip:       "127.0.0.1",
		Port:     port,
		Name:     "client1",
		Hash:     createClientHash("127.0.0.1", port, "client1"),
	}

	me = cl
	// register me as the first client
	cls.addClient(cl)
	// fetch a list of existing Clients
	cls.syncClients()
	// register me at all other Clients
	cls.greetClients()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(port uint16) {
	p := fmt.Sprintf("%d", port)
	fmt.Println("Starting server")
	fmt.Printf("Running on Port %s\n", p)
	log.Fatal(http.ListenAndServe(":"+p, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.index).Methods("GET")
	// transactions
	a.Router.HandleFunc("/transaction", a.newTransaction).Methods("POST")
	// blocks
	a.Router.HandleFunc("/block", a.connectClient).Methods("POST")
	// mining and chaining
	a.Router.HandleFunc("/mine", a.mine).Methods("GET")
	a.Router.HandleFunc("/chain", a.chain).Methods("GET")
	a.Router.HandleFunc("/validate", a.validate).Methods("GET")
	// Clients
	a.Router.HandleFunc("/client", a.connectClient).Methods("POST")
	a.Router.HandleFunc("/client", a.getClients).Methods("GET")
}

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Hello world")
}

// connectClient Connect a Client to the network which is represented
// in the Clients.list The postdata should consist of a standard Client
func (a *App) connectClient(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newCl Client
	err := decoder.Decode(&newCl)
	if err != nil {
		messenger("Could not decode postdata of new client")
		respondWithError(w, http.StatusBadRequest, "invalid json")
		panic(err)
	}

	newCl.Hash = createClientHash(newCl.Ip, newCl.Port, newCl.Name)
	// register the client
	cls.addClient(newCl)

	resp := map[string]interface{}{"Client": newCl, "total": cls.num()}
	respondWithJSON(w, http.StatusOK, resp)
}

// getClients response is the list of Clients
func (a *App) getClients(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{"list": cls.List, "length": len(cls.List)}
	respondWithJSON(w, http.StatusOK, resp)
}

// newTransaction adds a transaction, which consists of:
// Sender string
// Recipient string
// Amount float64
func (a *App) newTransaction(w http.ResponseWriter, r *http.Request) {

	var tr Transaction

	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Transaction")
	}

	bc.newTransaction(tr)

	respondWithJSON(w, http.StatusOK, "Transaction added")
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

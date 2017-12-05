package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/grrrben/golog"
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

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	// add the Client to the stack
	cls = initClients() // a pointer to the Clients struct
	cl := Client{
		Protocol: "http://",
		Hostname: name,
		Port:     clientPort,
		Name:     *clientName,
		Hash:     createClientHash(),
	}

	me = cl
	// register me as the first client
	cls.addClient(cl)
	// fetch a list of existing Clients
	cls.syncClients()
	// register me at all other Clients
	cls.greetClients()

	bc = initBlockchain()
	bc.getCurrentTransactions()
	golog.Info("Starting with a base blockchain:")
	golog.Infof("Blockchain:\n %v\n", bc)
	golog.Flush()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run() {
	p := fmt.Sprintf("%d", clientPort)
	fmt.Println("Starting server")
	fmt.Printf("Running on Port %s\n", p)
	log.Fatal(http.ListenAndServe(":"+p, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.index).Methods("GET")
	// transactions
	a.Router.HandleFunc("/transaction", a.newTransaction).Methods("POST")
	a.Router.HandleFunc("/transaction/distributed", a.distributedTransaction).Methods("POST")
	a.Router.HandleFunc("/transactions/{hash}", a.transactions).Methods("GET")
	a.Router.HandleFunc("/transactions", a.currentTransactions).Methods("GET")
	// wallet
	a.Router.HandleFunc("/wallet/{hash}", a.wallet).Methods("GET")
	// blocks
	a.Router.HandleFunc("/block", a.lastblock).Methods("GET")
	a.Router.HandleFunc("/block/{hash}", a.block).Methods("GET")
	a.Router.HandleFunc("/block/index/{index}", a.blockByIndex).Methods("GET")
	a.Router.HandleFunc("/block/distributed", a.distributedBlock).Methods("POST")
	// mining and chaining
	a.Router.HandleFunc("/mine", a.mine).Methods("GET")
	a.Router.HandleFunc("/chain", a.chain).Methods("GET")
	a.Router.HandleFunc("/validate", a.validate).Methods("GET")
	a.Router.HandleFunc("/resolve", a.resolve).Methods("GET")
	a.Router.HandleFunc("/status", a.chainStatus).Methods("GET")
	// Clients
	a.Router.HandleFunc("/client", a.connectClient).Methods("POST")
	a.Router.HandleFunc("/client", a.getClients).Methods("GET")
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

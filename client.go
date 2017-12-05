package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grrrben/golog"
	"net/http"
)

// this is me, a client
var me Client

type Client struct {
	Hostname string `json:"hostname"`
	Protocol string `json:"protocol"`
	Port     uint16 `json:"port"`
	Name     string `json:"name"`
	Hash     string `json:"hash"`
}

type Clients struct {
	List []Client
}

func initClients() *Clients {
	cls := &Clients{}
	return cls
}

// addClient Add a new Client to the list.
// A Client can only be added a single time, the list is unique.
// return bool true on success.
func (cls *Clients) addClient(cl Client) bool {
	for _, c := range cls.List {
		if c.Hash == cl.Hash {
			return false
		}
	}
	cls.List = append(cls.List, cl)
	golog.Infof("Client added. Clients: %d\n", cls.num())
	return true
}

// syncClients contacts other Clients to fetch a full list of Clients
// todo; how do I know which nodes are currently in the network
func (cls *Clients) syncClients() bool {
	// if I am the Mother node, ignore this
	if me.Port == 8000 {
		return true
	}
	// for now, just use a main parent node
	url := "http://localhost:8000/client"

	var externalCls Clients

	resp, err := http.Get(url)
	if err != nil {
		golog.Warningf("Could not get list of Clients on url: %s", url)
		return false
	}
	defer resp.Body.Close()
	decodingErr := json.NewDecoder(resp.Body).Decode(&externalCls)
	if decodingErr != nil {
		golog.Warningf("Could not decode JSON of list of Clients\n")
		return false
	}

	golog.Infof("externalCls:\n%v\n", externalCls)

	// just try to add all clients
	i := 0
	for _, c := range externalCls.List {
		success := cls.addClient(c)
		if success == true {
			i++
		}
	}
	golog.Infof("%d external Client(s) added\n", i)
	return true
}

// greetClients contacts other Clients to add this client to their list of known Clients
func (cls *Clients) greetClients() bool {
	for _, cl := range cls.List {
		if cl == me {
			// no need to register myself
			continue
		}
		go greet(cl)
	}
	return true
}

// greet makes a call to a client cl to make this node known within the network.
func greet(cl Client) {
	// POST to /client
	url := fmt.Sprintf("%s/client", cls.getAddress(cl))
	payload, err := json.Marshal(me)
	if err != nil {
		golog.Warning("Could not marshall client: Me")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		golog.Warningf("Request setup error: %s", err)
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Warningf("POST request error: %s", err)
		// I don't want to panic here, but it might be a good idea to
		// remove the client from the list
	}
	defer resp.Body.Close()
}

// announceMinedBlocks tells all clients in the network about the newly mined block.
// it gives the new block to the clients who can add it to their chain.
func (cls *Clients) announceMinedBlocks(bl Block) {
	for _, cl := range cls.List {
		if cl == me {
			continue // no need to brag
		}
		go announceMinedBlock(cl, bl)
	}
}

// distributeTransaction tells all clients in the network about the new Transaction.
func (cls *Clients) distributeTransaction(tr Transaction) {
	golog.Infof("Announcing transaction to %d clients", len(cls.List))
	for _, cl := range cls.List {
		if cl == me {
			continue // no need to brag
		}
		golog.Info("Announcing transaction")
		go announceTransaction(cl, tr)
	}
}

// num returns an int which represents the number of connected clients.
func (cls *Clients) num() int {
	return len(cls.List)
}

// getAddress returns (URI) address of a client.
func (cls *Clients) getAddress(cl Client) string {
	return fmt.Sprintf("%s%s:%d", cl.Protocol, cl.Hostname, cl.Port)
}

// createClientHash Creates a wallet and returns the hash of the new wallet.
// If a clients mines a block, the incentive is sent to this wallet address
func createClientHash() string {
	wallet := createWallet()
	return wallet.hash
}

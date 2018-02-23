package main

import (
	"encoding/json"
	"github.com/grrrben/golog"
	"net/http"
)

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
	cl.createWallet()
	for _, c := range cls.List {
		if c.getAddress() == cl.getAddress() {
			golog.Warningf("Client already known: %s", c.getAddress())
			return false
		}
	}
	cls.List = append(cls.List, cl)
	golog.Infof("Client added (%s). Clients: %d\n", cl.getAddress(), cls.num())
	return true
}

// syncClients contacts other Clients to fetch a full list of Clients
// todo; how do I know which nodes are currently in the network
func (cls *Clients) syncClients() bool {
	// if I am the only node, ignore this
	if me.Port == 8000 {
		return true
	}
	// for now, just use the main parent node as an oracle.
	url := "http://localhost:8000/client"

	var externalCls Clients

	resp, err := http.Get(url)
	if err != nil {
		golog.Warningf("Could not get list of Clients on url: %s", url)
		return false
	}

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

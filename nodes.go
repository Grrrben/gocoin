package main

import (
	"encoding/json"
	"net/http"

	"github.com/grrrben/golog"
)

type Nodes struct {
	List []Node
}

func initNodes() *Nodes {
	cls := &Nodes{}
	return cls
}

// addNode Add a new Node to the list.
// A N ode can only be added a single time, the list is unique.
// return bool true on success.
func (cls *Nodes) addNode(cl Node) bool {
	cl.createWallet()
	for _, c := range cls.List {
		if c.getAddress() == cl.getAddress() {
			golog.Warningf("Node already known: %s", c.getAddress())
			return false
		}
	}
	cls.List = append(cls.List, cl)
	golog.Infof("Node added (%s). Nodes: %d\n", cl.getAddress(), cls.num())
	return true
}

// syncNodes contacts other Nodes to fetch a full list of Nodes
// todo; how do I know which nodes are currently in the network
func (cls *Nodes) syncNodes() bool {
	// if I am the only node, ignore this
	if me.Port == 8000 {
		return true
	}
	// for now, just use the main parent node as an oracle.
	url := "http://localhost:8000/node"

	var externalCls Nodes

	resp, err := http.Get(url)
	if err != nil {
		golog.Warningf("Could not get list of Nodes on url: %s", url)
		return false
	}

	decodingErr := json.NewDecoder(resp.Body).Decode(&externalCls)
	if decodingErr != nil {
		golog.Warningf("Could not decode JSON of list of Nodes\n")
		return false
	}

	golog.Infof("externalCls:\n%v\n", externalCls)

	// just try to add all nodes
	i := 0
	for _, c := range externalCls.List {
		success := cls.addNode(c)
		if success == true {
			i++
		}
	}
	golog.Infof("%d external Node(s) added\n", i)
	return true
}

// greetNodes contacts other Nodes to add this node to their list of known Nodes
func (cls *Nodes) greetNodes() bool {
	for _, cl := range cls.List {
		if cl == me {
			// no need to register myself
			continue
		}
		go greet(cl)
	}
	return true
}

// announceMinedBlocks tells all nodes in the network about the newly mined block.
// it gives the new block to the nodes who can add it to their chain.
func (cls *Nodes) announceMinedBlocks(bl Block) {
	for _, cl := range cls.List {
		if cl == me {
			continue // no need to brag
		}
		go announceMinedBlock(cl, bl)
	}
}

// distributeTransaction tells all nodes in the network about the new Transaction.
func (cls *Nodes) distributeTransaction(tr Transaction) {
	golog.Infof("Announcing transaction to %d nodes", len(cls.List))
	for _, cl := range cls.List {
		if cl == me {
			continue // no need to brag
		}
		golog.Info("Announcing transaction")
		go announceTransaction(cl, tr)
	}
}

// num returns an int which represents the number of connected nodes.
func (cls *Nodes) num() int {
	return len(cls.List)
}

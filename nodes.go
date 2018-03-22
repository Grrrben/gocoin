package main

import (
	"encoding/json"
	"net/http"

	"github.com/grrrben/glog"
)

type Nodes struct {
	List []Node
}

func initNodes() *Nodes {
	cls := &Nodes{}
	return cls
}

// addNode Add a new Node to the list.
// A Node can only be added a single time, the list is unique.
// return bool true on success.
func (nodes *Nodes) addNode(newNode *Node) bool {
	newNode.createWallet()
	for _, n := range nodes.List {
		if n.getAddress() == newNode.getAddress() {
			glog.Warningf("Node already known: %s", n.getAddress())
			return false
		}
	}
	nodes.List = append(nodes.List, *newNode)
	glog.Infof("Node added (%s). Nodes: %d", newNode.getAddress(), nodes.num())
	return true
}

// syncNodes contacts other Nodes to fetch a full list of Nodes
// todo; how do I know which nodes are currently in the network
func (nodes *Nodes) syncNodes() bool {
	// if I am the only node, ignore this
	if me.Port == 8000 {
		return true
	}
	// for now, just use the main parent node as an oracle.
	url := "http://localhost:8000/node"

	var externalNodes Nodes

	resp, err := http.Get(url)
	if err != nil {
		glog.Warningf("Could not get list of Nodes on url: %s", url)
		return false
	}

	decodingErr := json.NewDecoder(resp.Body).Decode(&externalNodes)
	if decodingErr != nil {
		glog.Warningf("Could not decode JSON of list of Nodes\n")
		return false
	}

	glog.Infof("externalCls:\n%v\n", externalNodes)

	// just try to add all nodes
	i := 0
	for _, n := range externalNodes.List {
		success := nodes.addNode(&n)
		if success == true {
			i++
		}
	}
	glog.Infof("%d external Node(s) added\n", i)
	return true
}

// greetNodes contacts other Nodes to add this node to their list of known Nodes
func (nodes *Nodes) greetNodes() bool {
	for _, node := range nodes.List {
		if node == me {
			// no need to register myself
			continue
		}
		go greet(node)
	}
	return true
}

// announceMinedBlocks tells all nodes in the network about the newly mined block.
// it gives the new block to the nodes who can add it to their chain.
func (nodes *Nodes) announceMinedBlocks(bl Block) {
	for _, node := range nodes.List {
		if node == me {
			continue // no need to brag
		}
		go announceMinedBlock(node, bl)
	}
}

// distributeTransaction tells all nodes in the network about the new Transaction.
func (nodes *Nodes) distributeTransaction(tr Transaction) {
	glog.Infof("Announcing transaction to %d nodes", len(nodes.List))
	for _, node := range nodes.List {
		if node == me {
			continue // no need to brag
		}
		glog.Info("Announcing transaction")
		go announceTransaction(node, tr)
	}
}

// num returns an int which represents the number of connected nodes.
func (nodes *Nodes) num() int {
	return len(nodes.List)
}

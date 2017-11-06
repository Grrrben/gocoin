package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/grrrben/golog"
	"net/http"
)

// this is me, a client
var me Client

// @todo Clients via IP
type Client struct {
	Ip       string
	Protocol string
	Port     uint16
	Name     string
	Hash     string // todo should be the wallet Hash?
}

type Clients struct {
	List []Client
}

type clientService interface {
	addClient(Client) bool
	num() int
	syncClients() bool
	greetClients() bool
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
	golog.Infof("Client body:\n%v\n", resp.Body)
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
func greet (cl Client) {
	// POST to /client
	url := fmt.Sprintf("%s%s:%d/client", cl.Protocol, cl.Ip, cl.Port)
	golog.Infof("client URL: %s\n", url)

	payload, err := json.Marshal(me)
	golog.Infof("Me: %v\n", me)
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
		// I don't want to panic here, but it could be a good idea to
		// remove the client from the list
	}
	defer resp.Body.Close()
}

// num returns an int which represents the number of connected clients.
func (cls *Clients) num() int {
	return len(cls.List)
}

// createClientHash
// todo; Should be wallet?
func createClientHash(ip string, port uint16, name string) string {
	id := fmt.Sprint("%s-%d-%s", ip, port, name)

	jsonId, errr := json.Marshal(id)
	if errr != nil {
		if debug {
			golog.Errorf("Error: %s", errr)
		}
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, jsonId)
	if err != nil {
		if debug {
			golog.Errorf("Could not compute Client Hash. Msg: %s", err)
		}
	}
	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())) // %x; base 16, with lower-case letters for a-f
}

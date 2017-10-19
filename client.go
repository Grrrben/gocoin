package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
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

// used to fetch the post data from the Client
type ClientPost struct {
	Name string
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
			messenger("Client already exist in the network")
			return false
		}
	}
	cls.List = append(cls.List, cl)
	messenger("Client added. Clients: %d", cls.num())
	return true
}

// syncClients contacts other Clients to fetch a full list of Clients
// todo; how do I know which nodes are currently in the network
func (cls *Clients) syncClients() bool {
	// for now, just use a main parent node
	url := "http://localhost:8000"

	var externalCls Clients

	resp, err := http.Get(url)
	if err != nil {
		messenger("Could not get list of Clients on url: %s", url)
		return false
	}
	defer resp.Body.Close()
	decodingErr := json.NewDecoder(resp.Body).Decode(&externalCls)
	if decodingErr != nil {
		messenger("Could not decode JSON of list of Clients")
		return false
	}

	// just try to add all clients
	i := 0
	for _, c := range externalCls.List {
		success := cls.addClient(c)
		if success == true {
			i++
		}
	}
	messenger("%d external Client(s) added", i)

	return true
}

// greetClients contacts other Clients to add this client to their list of known Clients
func (cls *Clients) greetClients() bool {
	for _, cl := range cls.List {
		if cl == me {
			// no need to register myself
			continue
		}
		// POST to /client
		url := fmt.Sprintf("%s%s%s/client", cl.Protocol, cl.Ip, cl.Port)
		messenger("client URL: %s", url)

		payload, err := json.Marshal(me)
		if err != nil {
			messenger("Could not marshall client: Me")
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			messenger("POST request error: %s", err)
			panic(err)
		}

		fmt.Printf("Resp %v", resp.Body)

		defer resp.Body.Close()

	}
	return true
}

func (cls *Clients) num() int {
	return len(cls.List)
}

func createClientHash(ip string, port uint16, name string) string {
	id := fmt.Sprint("%s-%d-%s", ip, port, name)

	jsonId, errr := json.Marshal(id)
	if errr != nil {
		if debug {
			fmt.Printf("Error: %s", errr)
		}
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, jsonId)
	if err != nil {
		if debug {
			fmt.Println("Could not compute Client Hash")
			fmt.Println(err)
		}
	}
	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())) // %x; base 16, with lower-case letters for a-f
}

package main

import (
	"flag"
	"github.com/grrrben/golog"
	"strconv"
)

var bc *Blockchain
var cls *Clients

func main() {
	prt := flag.String("p", "8000", "Port on which the app will run, defaults to 8000")
	clientName:= flag.String("name", "0", "Set a name for the client")
	flag.Parse()

	golog.SetLogDir("/home/grrrben/go/src/blockchain/log")

	u, err := strconv.ParseUint(*prt, 10, 16) // always gives an uint64...
	if err != nil {
		golog.Errorf("Unable to cast Prt to uint: %s", err)
	}
	// different Clients can have different ports,
	// used to connect multiple Clients in debug.
	clientPortNr := uint16(u)

	a := App{}
	a.Initialize(clientPortNr, *clientName)
	a.Run(clientPortNr)
}

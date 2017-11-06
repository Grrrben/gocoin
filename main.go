package main

import (
	"flag"
	"github.com/grrrben/golog"
	"strconv"
)

var bc *Blockchain
var cls *Clients
var debug bool

func main() {
	prt := flag.String("p", "8000", "Port on which the app will run, defaults to 8000")
	verbose := flag.String("verbose", "0", "Verbose, show debug messages when set (1)")
	flag.Parse()

	if *verbose == "1" {
		debug = true
	}

	golog.SetLogDir("/home/grrrben/go/src/blockchain/log")
	golog.Info("golog info line")
	golog.Warning("golog warning line")
	golog.Flush()

	u, err := strconv.ParseUint(*prt, 10, 16) // always gives an uint64...
	if err != nil {
		golog.Errorf("Unable to cast Prt to uint: %s", err)
	}
	// different Clients can have different ports,
	// used to connect multiple Clients in debug.
	clientPortNr := uint16(u)

	a := App{}
	a.Initialize(clientPortNr)
	a.Run(clientPortNr)
}

package main

import (
	"flag"
	"fmt"
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

	u, err := strconv.ParseUint(*prt, 10, 16) // always gives an uint64...
	if err != nil {
		fmt.Println("Unable to cast Prt to uint")
		fmt.Println(err)
	}
	// different Clients can have different ports,
	// used to connect multiple Clients in debug.
	clientPortNr := uint16(u)

	a := App{}
	a.Initialize(clientPortNr)
	a.Run(clientPortNr)
}

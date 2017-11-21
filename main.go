package main

import (
	"flag"
	"github.com/grrrben/golog"
	"strconv"
	"path/filepath"
	"os"
	"fmt"
)

var bc *Blockchain
var cls *Clients

var clientPort uint16
var clientName *string

func main() {
	prt := flag.String("p", "8000", "Port on which the app will run, defaults to 8000")
	clientName = flag.String("name", "0", "Set a name for the client")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		golog.Fatalf("Could not set a logdir. Msg %s", err)
	}

	golog.SetLogDir(fmt.Sprintf("%s/log", dir))

	u, err := strconv.ParseUint(*prt, 10, 16) // always gives an uint64...
	if err != nil {
		golog.Errorf("Unable to cast Prt to uint: %s", err)
	}
	// different Clients can have different ports,
	// used to connect multiple Clients in debug.
	clientPort = uint16(u)

	a := App{}
	a.Initialize()
	a.Run()
}

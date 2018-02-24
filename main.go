package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/grrrben/golog"
)

var bc *Blockchain
var cls *Nodes

var nodePort uint16
var nodeName *string

const zerohash = "0000000000000000000000000000000000000000000000000000000000000000"

func main() {
	prt := flag.String("p", "8000", "Port on which the app will run, defaults to 8000")
	nodeName = flag.String("name", "0", "Set a name for the node")
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
	// different Nodes can have different ports,
	// used to connect multiple Nodes in debug.
	nodePort = uint16(u)

	a := App{}
	a.Initialize()
	a.Run()
}

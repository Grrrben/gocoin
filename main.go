package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"log"

	"github.com/grrrben/glog"
)

var bc *Blockchain
var nodes *Nodes

var nodePort uint16
var nodeName *string

const zerohash = "0000000000000000000000000000000000000000000000000000000000000000"

func main() {
	prt := flag.String("p", "8000", "Port on which the app will run, defaults to 8000")
	nodeName = flag.String("name", "Node_X", "Set a name for the node")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Could not set a logdir. Msg %s", err.Error())
	}

	glog.SetLogFile(fmt.Sprintf("%s/log/blockchain.log", dir))
	glog.SetLogLevel(glog.Log_level_warning)

	u, err := strconv.ParseUint(*prt, 10, 16) // always gives an uint64...
	if err != nil {
		glog.Errorf("Unable to cast Prt to uint: %s", err.Error())
	}
	// different Nodes can have different ports,
	// used to connect multiple Nodes in debug.
	nodePort = uint16(u)

	a := App{}
	a.Initialize()
	a.Run()
}

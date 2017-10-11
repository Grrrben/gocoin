package main

import "flag"

var bc *Blockchain
var debug bool

func main() {
	//serverPort := flag.String("p", "8080", "Port on which the server will run, defaults to 8080")
	verbose := flag.String("verbose", "0", "Verbose, show debug messages when set (1)")
	flag.Parse()

	if *verbose == "1" {
		debug = true
	}

	a := App{}
	a.Initialize()
	a.Run()
}

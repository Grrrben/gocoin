package main

import "fmt"

func messenger (message string, vars ...interface{}) {
	// todo; this works but it should differentiate between debug, warning, error
	// todo; add logging.
	if debug {
		if len(vars) > 0 {
			fmt.Printf(message, vars...)
		} else {
			fmt.Println(message)
		}
	}
}

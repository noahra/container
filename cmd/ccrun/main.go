package main

import (
	"flag"
)

func main() {
	rFlag := flag.Bool("r", false, "Enable r option")
	flag.Parse()
	args := flag.Args()
	if *rFlag {
		setHostname("container123")
		executeCommand(args)
	} else {
		createUtsNameSpace(args)
	}
}

package main

import (
	"flag"
	"fmt"
)

func main() {
	rFlag := flag.Bool("r", false, "Enable r option")
	flag.Parse()
	args := flag.Args()
	if *rFlag {
		setHostname("container123")
		setChroot("alpine_fs")
		err := executeCommand(args)
		if err != nil {
			fmt.Printf("error occured when executing command: %s", err)
		}
	} else {
		createUtsNameSpace(args)
	}
}

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
		err := unShareMount()
		if err != nil {
			fmt.Printf("error occured when unsharing mount: %s", err)
		}
		err = setHostname("container123")
		if err != nil {
			fmt.Printf("error occured when setting hostname (creating uts ns): %s", err)
		}
		err = setChroot("alpine_fs")
		if err != nil {
			fmt.Printf("error occured when executing chroot: %s", err)
		}
		err = mountProc()
		if err != nil {
			fmt.Printf("error occured when mounting /proc: %s", err)
		}

		err = executeCommand(args)
		if err != nil {
			fmt.Printf("error occured when executing command: %s", err)
		}
	} else {
		createNameSpaces(args)
	}
}

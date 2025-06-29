package main

import (
	"flag"
	"fmt"
)

func main() {
	const HOSTNAME = "container123"
	rFlag := flag.Bool("r", false, "Re-execute as container init process")
	flag.Parse()
	args := flag.Args()
	if *rFlag {
		err := setHostname(HOSTNAME)
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
		cmd := createNameSpaces(args)
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
		}
		cg(cmd.Process.Pid, HOSTNAME)
		cmd.Wait()

		cleanupCgroups(HOSTNAME)
	}
}

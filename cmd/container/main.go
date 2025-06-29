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

	if !*rFlag {
		cmd := createNameSpaces(args)
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			return
		}
		cg(cmd.Process.Pid, HOSTNAME)
		cmd.Wait()
		cleanupCgroups(HOSTNAME)
		return
	}

	// Container init process path
	if err := setHostname(HOSTNAME); err != nil {
		fmt.Printf("error occurred when setting hostname: %s\n", err)
		return
	}
	if err := setChroot("alpine_fs"); err != nil {
		fmt.Printf("error occurred when executing chroot: %s\n", err)
		return
	}
	if err := mountProc(); err != nil {
		fmt.Printf("error occurred when mounting /proc: %s\n", err)
		return
	}
	if err := executeCommand(args); err != nil {
		fmt.Printf("error occurred when executing command: %s\n", err)
	}
}

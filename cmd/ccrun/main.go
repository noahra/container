package main

import (
	"flag"
	"fmt"
)

func main() {
	const HOSTNAME = "container123"
	rFlag := flag.Bool("r", false, "Enable r option")
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
		fmt.Printf("DEBUG: Creating container process...\n")
		cmd := createNameSpaces(args)

		fmt.Printf("DEBUG: Waiting for container to finish...\n")
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("DEBUG: Creating cgroup from host for PID %d\n", cmd.Process.Pid)
		cg(cmd.Process.Pid, HOSTNAME)

		cmd.Wait()
		fmt.Printf("DEBUG: Cleaning up cgroups...\n")

		cleanupCgroups(HOSTNAME)
	}
}

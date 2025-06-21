package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	const HOSTNAME = "container123"
	rFlag := flag.Bool("r", false, "Enable r option")
	pidFlag := flag.String("pid", "", "PID to add to cgroup")
	flag.Parse()
	args := flag.Args()
	if *rFlag {
		if *pidFlag == "self" {
			pid := os.Getpid()
			cg(pid, HOSTNAME)
		} else if *pidFlag != "" {
			pid, err := strconv.Atoi(*pidFlag)
			if err != nil {
				fmt.Printf("Invalid PID: %s\n", err)
				return
			}
			cg(pid, HOSTNAME)
		}

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
		cmd.Wait()
		cleanupCgroups(HOSTNAME)

	}
}

func cg(pid int, hostname string) {
	// Create single cgroup directory (all controllers in one place)
	cgroupPath := filepath.Join("/sys/fs/cgroup", hostname)
	err := os.MkdirAll(cgroupPath, 0755)
	if err != nil {
		fmt.Printf("Error creating cgroup: %s\n", err)
		return
	}

	// CPU limits
	cpuQuotaFile := filepath.Join(cgroupPath, "cpu.cfs_quota_us")
	err = os.WriteFile(cpuQuotaFile, []byte("50000"), 0644)
	if err != nil {
		fmt.Printf("Error setting CPU quota: %s\n", err)
		return
	}

	cpuPeriodFile := filepath.Join(cgroupPath, "cpu.cfs_period_us")
	err = os.WriteFile(cpuPeriodFile, []byte("100000"), 0644)
	if err != nil {
		fmt.Printf("Error setting CPU period: %s\n", err)
		return
	}

	// PID limits - limit to 10 processes/threads
	pidsMaxFile := filepath.Join(cgroupPath, "pids.max")
	err = os.WriteFile(pidsMaxFile, []byte("10"), 0644)
	if err != nil {
		fmt.Printf("Error setting PID limit: %s\n", err)
		return
	}

	// Add PID to cgroup
	procsFile := filepath.Join(cgroupPath, "cgroup.procs")
	err = os.WriteFile(procsFile, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		fmt.Printf("Error adding PID to cgroup: %s\n", err)
		return
	}
}

func cleanupCgroups(hostname string) {
	// Remove single cgroup folder
	cgroupPath := filepath.Join("/sys/fs/cgroup", hostname)
	err := os.RemoveAll(cgroupPath)
	if err != nil {
		fmt.Printf("Error removing cgroup: %s\n", err)
	}
}

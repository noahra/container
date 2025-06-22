package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
func cg(pid int, hostname string) {
	// CPU Controller
	cpuCgroupPath := filepath.Join("/sys/fs/cgroup/cpu", hostname)
	err := os.MkdirAll(cpuCgroupPath, 0755)
	if err != nil {
		fmt.Printf("Error creating CPU cgroup: %s\n", err)
		return
	}

	// Set CPU limits
	cpuQuotaFile := filepath.Join(cpuCgroupPath, "cpu.cfs_quota_us")
	err = os.WriteFile(cpuQuotaFile, []byte("50000"), 0644)
	if err != nil {
		fmt.Printf("Error setting CPU quota: %s\n", err)
		return
	}

	cpuPeriodFile := filepath.Join(cpuCgroupPath, "cpu.cfs_period_us")
	err = os.WriteFile(cpuPeriodFile, []byte("100000"), 0644)
	if err != nil {
		fmt.Printf("Error setting CPU period: %s\n", err)
		return
	}
	// Set PID limits
	pidsMaxFile := filepath.Join(cpuCgroupPath, "pids.max")
	err = os.WriteFile(pidsMaxFile, []byte("20"), 0644)
	if err != nil {
		fmt.Printf("Error setting PID limit: %s\n", err)
		return
	}

	// Add PID to CPU cgroup
	cpuProcsFile := filepath.Join(cpuCgroupPath, "cgroup.procs")
	err = os.WriteFile(cpuProcsFile, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		fmt.Printf("Error adding PID to CPU cgroup: %s\n", err)
		return
	}

}

func cleanupCgroups(hostname string) {
	// Clean up CPU cgroup
	cpuCgroupPath := filepath.Join("/sys/fs/cgroup/cpu", hostname)
	err := os.RemoveAll(cpuCgroupPath)
	if err != nil {
		fmt.Printf("Error removing CPU cgroup: %s\n", err)
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func cg(pid int, hostname string) {

	rootCgroupPath := "/sys/fs/cgroup"
	hostnameCgroupPath := filepath.Join(rootCgroupPath, hostname)

	err := os.MkdirAll(hostnameCgroupPath, 0755)
	if err != nil {
		fmt.Printf("Error creating container's cgroup: %s\n", err)
		return
	}

	cpuProcsFile := filepath.Join(hostnameCgroupPath, "cgroup.procs")
	err = os.WriteFile(cpuProcsFile, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		fmt.Printf("Error adding PID to CPU cgroup: %s\n", err)
		return
	}

	cpuMaxFile := filepath.Join(hostnameCgroupPath, "cpu.max")
	// "50000 100000" means 50% of one CPU (50ms out of every 100ms period)
	err = os.WriteFile(cpuMaxFile, []byte("50000 100000"), 0644)
	if err != nil {
		fmt.Printf("Error setting CPU max: %s\n", err)
		return
	}

}

func cleanupCgroups(hostname string) {
	cpuCgroupPath := filepath.Join("/sys/fs/cgroup", hostname)
	err := os.RemoveAll(cpuCgroupPath)
	if err != nil {
		fmt.Printf("Error removing CPU cgroup: %s\n", err)
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func cg(pid int, hostname string) {

	subtreeControlTreePath := "/sys/fs/cgroup"

	subtreeControlFile := filepath.Join(subtreeControlTreePath, "cgroup.subtree_control")

	current, err := os.ReadFile(subtreeControlFile)
	if err != nil {
		fmt.Printf("Error reading current state: %s\n", err)
		return
	}
	// Parse and add new controllers (more complex logic needed here)
	newContent := string(current) + " +pids"
	err = os.WriteFile(subtreeControlFile, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error enabling controllers: %s\n", err)
		return
	}

	hostnameCgroupPath := filepath.Join("/sys/fs/cgroup/", hostname)
	err = os.MkdirAll(hostnameCgroupPath, 0755)
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

	pidsMaxFile := filepath.Join(hostnameCgroupPath, "pids.max")
	err = os.WriteFile(pidsMaxFile, []byte("20"), 0644)
	if err != nil {
		fmt.Printf("Error setting PIDs limit: %s\n", err)
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

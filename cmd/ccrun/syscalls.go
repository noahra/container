package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func createUtsNameSpace(args []string) {
	cmd := exec.Cmd{
		Path:   "/proc/self/exe", // Just the executable path as a string
		Args:   append([]string{"/proc/self/exe", "-r"}, args[1:]...),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS,
		},
	}
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func setHostname(hostname string) {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Failed to set hostname: %v\n", err)
		return
	}
	fmt.Printf("Hostname set to: %s\n", hostname)
}

func setChroot(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to set absolute path with path: %s: %w", path, err)
	}
	if err = os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change to directory %s: %v", path, err)
	}

	if err = syscall.Chroot(absolutePath); err != nil {
		return fmt.Errorf("failed to chroot to %s: %v", path, err)
	}

	return nil
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func createNameSpaces(args []string) exec.Cmd {
	cmd := exec.Cmd{
		Path:   "/proc/self/exe",
		Args:   append([]string{"/proc/self/exe", "-r"}, args[1:]...),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
			UidMappings: []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: os.Getuid(), Size: 1},
			},
			GidMappings: []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: os.Getgid(), Size: 1},
			},
		},
	}
	return cmd
}

func setHostname(hostname string) error {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Failed to set hostname: %v\n", err)
		return err
	}
	fmt.Printf("Hostname set to: %s\n", hostname)
	return nil
}

func setChroot(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to set absolute path with path: %s: %w", path, err)
	}

	if err := os.Chdir(absolutePath); err != nil {
		return fmt.Errorf("failed to change to directory %s: %w", path, err)
	}

	if err = syscall.Chroot(absolutePath); err != nil {
		return fmt.Errorf("failed to chroot to %s: %w", path, err)
	}

	return nil
}
func mountProc() error {
	const proc = "proc"
	absolutePathToProcFolder, err := filepath.Abs(proc)
	if err != nil {
		return fmt.Errorf("failed to parse absolute path to proc folder: %w", err)
	}
	if err := syscall.Mount(proc, absolutePathToProcFolder, proc, 0, ""); err != nil {
		return fmt.Errorf("failed to mount /proc: %w", err)
	}
	fmt.Printf("Proc mounted")
	return nil
}

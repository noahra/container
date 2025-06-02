package main

import (
	"fmt"
	"os"
	"os/exec"
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

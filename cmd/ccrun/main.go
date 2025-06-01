package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	argsWithoutProg := os.Args[1:]

	exec.Command("/bin/bash")
	cmd := exec.Cmd{
		Path:   "/bin/bash",
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

	if len(argsWithoutProg) >= 2 && argsWithoutProg[0] == "run" {

		// clone()
		// call command inside namespace
		executeCommand(argsWithoutProg)
	}
}

func executeCommand(args []string) {
	cmd := exec.Command(args[1], args[2:]...)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}

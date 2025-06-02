package main

import (
	"fmt"
	"os/exec"
)

func executeCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))
}

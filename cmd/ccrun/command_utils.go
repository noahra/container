package main

import (
	"fmt"
	"os/exec"
)

func executeCommand(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Print(string(output))
	return nil
}

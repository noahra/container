package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) >= 2 && argsWithoutProg[0] == "run" {
		cmd := exec.Command(argsWithoutProg[1], argsWithoutProg[2:]...)
		output, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Print(string(output))
	}
}

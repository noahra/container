package main

import (
	"os"
	"os/exec"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) >= 2 {
		cmd := exec.Command(argsWithoutProg[0], argsWithoutProg[:1])
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

}

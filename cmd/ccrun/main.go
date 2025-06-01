package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	rFlag := flag.Bool("r", false, "Enable r option")
	flag.Parse()
	args := flag.Args()
	fmt.Printf("Remaining args: %v\n", args)
	if *rFlag {
		println("checkpoint2")
		setHostname("container123")
		executeCommand(args[1:])
	} else {
		println("checkpoint1")
		createUtsNameSpace()
		commands := append([]string{"./main", "-r"}, args[1:]...)
		fmt.Println(commands)
		executeCommand(commands)
		fmt.Println("hmm")
	}
}

func createUtsNameSpace() {
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
}

func setHostname(hostname string) {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Failed to set hostname: %v\n", err)
		return
	}
	fmt.Printf("Hostname set to: %s\n", hostname)
}

func executeCommand(args []string) {
	fmt.Println("args: ", args)
	if len(args) >= 2 {
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Print(string(output))
	}
}

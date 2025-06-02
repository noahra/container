package main

import (
	"fmt"
	"syscall"
)

func setHostname(hostname string) {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Failed to set hostname: %v\n", err)
		return
	}
	fmt.Printf("Hostname set to: %s\n", hostname)
}

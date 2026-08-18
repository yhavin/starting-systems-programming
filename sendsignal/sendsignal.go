package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		fatal(fmt.Errorf("usage: %s <pid> <signal>", os.Args[0]))
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fatal(err)
	}
	signal, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fatal(err)
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_KILL,
		uintptr(pid),
		uintptr(signal),
		0,
	)
	if errno != 0 {
		fatal(errno)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
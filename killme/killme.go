package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	fmt.Println(pid)
	for {
		time.Sleep(15 * time.Second)
		fmt.Fprintln(os.Stderr, "still alive")
	}
}

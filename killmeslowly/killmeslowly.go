package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println(os.Getpid())

	ch := make(chan os.Signal, 5)
	signal.Notify(ch, syscall.SIGINT)

	fmt.Fprintf(os.Stderr, "waiting for SIGINT\n")
	remaining := 5
	for range ch {
		remaining --
		if remaining == 0 {
			fmt.Fprintf(os.Stderr, "exit\n")
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "got SIGINT: %d more to exit\n", remaining)
	}
}
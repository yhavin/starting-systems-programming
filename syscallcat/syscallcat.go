package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
	
	"starting-systems-programming/utils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	path := []byte(os.Args[1])
	path = append(path, 0)
	ptr := unsafe.Pointer(&path[0])
	const MODE = syscall.O_RDONLY
	const FLAGS = 0

	fileDescriptor, _, err := syscall.Syscall(
		syscall.SYS_OPEN,
		uintptr(unsafe.Pointer(ptr)),
		MODE,
		FLAGS,
	)
	if err != 0 {
		utils.Fatalf("open: %v\n", err)
	}

	var buf [1024]byte

READ:
	for {
		ptr := &buf[0]
		n, _, readErr := syscall.Syscall(
			syscall.SYS_READ,
			fileDescriptor,
			uintptr(unsafe.Pointer(ptr)),
			uintptr(len(buf)),
		)

		const FD_STDOUT = 1

		for offset := uintptr(0); offset < n; {
			ptr := &buf[offset]
			m, _, writeErr := syscall.Syscall(
				syscall.SYS_WRITE,
				FD_STDOUT,
				uintptr(unsafe.Pointer(ptr)),
				n,
			)
			if m == n {
				continue READ
			}
			if writeErr != 0 {
				utils.Fatalf("write: %v\n", writeErr)
			}
			offset += m
		}

		if readErr != 0 {
			utils.Fatalf("read: %v\n", readErr)
		}

		if n == 0 {
			break READ
		}
	}

	syscall.Syscall(syscall.SYS_FSYNC, 1, 0, 0)

	syscall.Syscall(syscall.SYS_CLOSE, fileDescriptor, 0, 0)

	syscall.Syscall(syscall.SYS_EXIT, 0, 0, 0)
}
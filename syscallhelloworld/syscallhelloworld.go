package main

import (
	"syscall"
	"unsafe"

	"starting-systems-programming/utils"
)

func main() {
	var b = []byte("hello, world!\n")
	n, _, errno := syscall.Syscall(
		syscall.SYS_WRITE,
		1,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
	)
	if errno != 0 {
		utils.Fatalf("write: %v\n", errno)
	}
	if n != uintptr(len(b)) {
		utils.Fatalf("write: wrote %d bytes, expected %d\n", n, len(b))
	}
}
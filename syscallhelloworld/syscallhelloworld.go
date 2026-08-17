package main

import (
	"fmt"
	"syscall"
	"unsafe"
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
		fatalf("write: %v\n", errno)
	}
	if n != uintptr(len(b)) {
		fatalf("write: wrote %d bytes, expected %d\n", n, len(b))
	}
}

func fatalf(format string, args ...interface{}) {
	buf := []byte(fmt.Sprintf(format, args...))
	syscall.Syscall(syscall.SYS_WRITE, uintptr(syscall.Stderr), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	syscall.Syscall(syscall.SYS_EXIT, 1, 0, 0)
}
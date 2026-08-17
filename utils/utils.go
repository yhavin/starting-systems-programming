package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

func Fatalf(format string, args ...interface{}) {
	buf := []byte(fmt.Sprintf(format, args...))
	syscall.Syscall(syscall.SYS_WRITE, uintptr(syscall.Stderr), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	syscall.Syscall(syscall.SYS_EXIT, 1, 0, 0)
}
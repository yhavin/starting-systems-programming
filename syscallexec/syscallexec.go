package main

import (
	"os"
	"syscall"
	"unsafe"

	"starting-systems-programming/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.Fatalf("usage: %s <command> [args...]\n", os.Args[0])
	}
	goargs := os.Args[1:]

	if len(os.Args[1]) == 0 || os.Args[1][0] != '/' {
		utils.Fatalf("error: %s is not an absolute path\n", os.Args[1])
	}
	exec(goargs, os.Environ())
}

func exec(args, env[]string) error {
	cargs := make([]unsafe.Pointer, len(args)+1)
	for i := range args {
		cargs[i] = utils.Cstr(args[i])
	}

	cenv := make([]unsafe.Pointer, len(env)+1)
	for i:= range env {
		cenv[i] = utils.Cstr(env[i])
	}

	path := utils.Cstr(args[0])

	_, _, err := syscall.Syscall(
		syscall.SYS_EXECVE,
		uintptr(path),
		uintptr(unsafe.Pointer(&cargs[0])),
		uintptr(unsafe.Pointer(&cenv[0])),
	)
	return err
}
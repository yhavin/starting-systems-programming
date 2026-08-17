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
		utils.Fatalf("usage: %s <command> [args...]\n", os.Args[0])
	}

	env := os.Environ()
	var path []string

	{
		var start int
		rawpath := getenv(env, "PATH")
		for i := range rawpath {
			if rawpath[i] != ':' {
				continue
			}
			if start == i {
				continue
			}
			path = append(path, rawpath[start:i])
			start = i + 1
		}
		if start < len(path) {
			path = append(path, rawpath[start:])
		}
	}

	for _, dir := range path {
		if path, err := exists(dir + "/" + os.Args[1]); err == nil && path {
			fmt.Println(dir + "/" + os.Args[1])
			os.Exit(0)
		}
	}

	utils.Fatalf("%s: command not found\n", os.Args[1])
}

func getenv(environ []string, key string) string {
	key += "="
	n := len(key)
	for i := range environ {
		if len(environ[i]) < len(key) {
			continue
		}
		if environ[i][:n] == key {
			return environ[i][n:]
		}
	}
	return ""
}

func exists(path string) (bool, error) {
	p := utils.Cstr(path)
	var statbuf [144]byte

	_, _, err := syscall.Syscall(
		syscall.SYS_STAT,
		uintptr(p),
		uintptr(unsafe.Pointer(&statbuf)),
		0,
	)

	switch err {
	case 0:
		return true, nil
	case syscall.ENOENT:
		return false, nil
	default:
		return false, err
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"starting-systems-programming/utils"
)

func main() {
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Fields(line)

		path, err := whiche(os.Environ(), args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%q: %v\n", args[0], err)
			continue
		}

		name := path

		args[0] = path

		status, err := call(args, os.Environ())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		}
		if status != 0 {
			fmt.Fprintf(os.Stderr, "%s: exit status %d\n", name, status)
		}
	}
}

func call(args []string, env []string) (status int, err error) {
	{
		if len(args) == 0 {
			return 0, errors.New("no command")
		}
		if args[0] == "" {
			return 0, errors.New("empty command")
		}
		if args[0][0] != '/' {
			return 0, errors.New("command must be an absolute path; try using lookupPath")
		}
	}

	pid, forkflag, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		return status, fmt.Errorf("syscall: fork: %v", errno)
	}

	const STATUS_FAILED_EXEC = 0xB01D

	if isChild := forkflag == 1; isChild {
		err := exec(args, env)
		fmt.Fprintf(os.Stderr, "syscall: execve: %v\n", err)
		os.Exit(STATUS_FAILED_EXEC)
	}

	{
		var waitstatus uint32
		pid, _, _ := syscall.Syscall6(
			syscall.SYS_WAIT4,
			pid,
			uintptr(unsafe.Pointer(&waitstatus)),
			0,
			0,
			0,
			0,
		)
		fmt.Fprintf(os.Stderr, "pid %d exited with status %d\n", pid, waitstatus)

		status = int(waitstatus >> 8)
		if status == STATUS_FAILED_EXEC {
			return status, errors.New("execve failed")
		}

		return status, nil
	}
}

func exec(args, env []string) error {
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

func whiche(env []string, command string) (string, error) {
	switch {
	case strings.Contains(command, ".."):
		return "", errors.New("no weird relative paths, please")
	case command == "":
		return "", errors.New("empty path")
	case command[0] == '/':
		return command, nil
	case command[0] == '.':
		wd, _ := os.Getwd()
		return lookupPath(command, wd)
	default:
		pathEnv := getenv(env, "PATH")
		dirs := strings.Split(pathEnv, ":")
		return lookupPath(command, dirs...)
	}
}

func getenv(environ []string, key string) string {
	key += "="
	n := len(key)
	for i:= range environ {
		if len(environ[i]) < len(key){
			continue
		}
		if environ[i][:n] == key {
			return environ[i][n:]
		}
	}
	return ""
}

func lookupPath(name string, dirs ...string) (string, error) {
	if strings.Contains(name, "..") {
		return "", fmt.Errorf("no weird relative paths, please: %q", name)
	}
	for i, dir := range dirs {
		path := dir + "/" + name
		if ok, err := exists(path); ok {
			return path, nil
		} else if err != nil {
			return "", fmt.Errorf("lookupPath: stat in dir #%d: %q: %w", i, path, err)
		}
	}
	return "", errors.New("not found in PATH")
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

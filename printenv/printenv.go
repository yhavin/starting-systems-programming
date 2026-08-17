package main

import (
	"fmt"
	"os"
)

func main() {
	env := os.Environ()
	var printed int
	keys := os.Args[1:]

	for _, key := range keys {
		val, ok := lookupenv(env, key)
		if ok {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", key, val)
			printed++
		}
	}

	if printed == len(keys) {
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "missing %d/%d environment variables\n", len(os.Args)-1-printed, len(os.Args)-1)

	for _, key := range keys {
		if _, ok := lookupenv(env, key); !ok {
			fmt.Fprintf(os.Stderr, "%s\n", key)
		}
	}
	os.Exit(1)
}

func lookupenv(env []string, key string) (string, bool) {
	for i := len(env)-1; i >= 0; i-- {
		e := env[i]
		if len(e) < len(key)+1 {
			continue
		}
		if e[:len(key)] != key {
			continue
		}
		if e[len(key)] != '=' {
			continue
		}
		return e[len(key)+1:], true
	}
	return "", false
}
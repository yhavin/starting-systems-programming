package main

import (
	"fmt"
	"os"
)

func main() {
	flags, positional, err := parseflags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parseflags: %v\n", err)
		os.Exit(1)
	}

	for name, value := range flags {
		fmt.Fprintf(os.Stdout, "flag %s=%s\n", name, value)
	}
	for i, arg := range positional {
		fmt.Fprintf(os.Stdout, "positional %d=%s\n", i, arg)
	}
}

func parseflags(args []string) (flags map[string]string, positional []string, err error) {
	flags = make(map[string]string)

	FLAGS:
		for len(args) > 0 {
			s := args[0]
			if len(s) <= 1 {
				break FLAGS
			}
			if s == "--" {
				args = args[1:]
				break FLAGS
			}
			if s[0] != '-' {
				break FLAGS
			}
			if s[1] == '-' {
				s = s[2:]
			} else {
				s = s[1:]
			}

			for i := range s {
				if s[i] == '=' {
					key, value := s[:i], s[i+1:]
					if _, ok := flags[key]; ok {
						return nil, nil, fmt.Errorf("flag -%s already set", key)
					}
					flags[key] = value
					args = args[1:]
					continue FLAGS
				}
			}

			if len(args) == 1 {
				return nil, nil, fmt.Errorf("flag -%s missing value", s)
			}

			key, value := s, args[1]
			if _, ok := flags[key]; ok {
				return nil, nil, fmt.Errorf("flag -%s already set", key)
			}
			flags[key] = value
			args = args[2:]
		}

		return flags, args, nil
}
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	var src io.Reader
	switch len(os.Args) {
	case 1:
		src= os.Stdin
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v", os.Args[1], err)
			os.Exit(1)
		}
		defer f.Close()
		src = f
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s [filename]", os.Args[0])
		os.Exit(1)
	}

	if err := unhexdump(os.Stdout, src); err != nil {
		fmt.Fprintf(os.Stderr, "unhexdump: %v", err)
		os.Exit(1)
	}
}

func unhexdump(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for i := 0; scanner.Scan(); i++ {
		b := scanner.Bytes()
		if len(b)&1 == 1 {
			return fmt.Errorf("odd number of hex digits as position %d", i)
		}

		for i := 0; i < len(b); i += 2 {
			high, ok := unhex(b[i])
			if !ok {
				return fmt.Errorf("bad hex %x '%c' at position %d", b[i], b[i], i)
			}
			low, ok := unhex(b[i+1])
			if !ok {
				return fmt.Errorf("bad hex %x '%c' at position %d", b[i+1], b[i+1], i+1)
			}

			if err := bw.WriteByte(high<<4 | low); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func unhex(b byte) (byte, bool) {
	switch {
	case '0' <= b && b <= '9':
		return b - '0', true
	case 'a' <= b && b <= 'f':
		return b - 'a' + 10, true
	case 'A' <= b && b <= 'F':
		return b - 'A' + 10, true
	default:
		return 0, false
	}
}
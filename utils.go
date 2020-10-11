package main

import (
	"io"
	"os"
)

func lineCount(s string) (n int, err error) {
	r, _ := os.Open(s)

	buf := make([]byte, 8192)

	defer r.Close()

	for {
		c, err := r.Read(buf)
		if err != nil {
			if err == io.EOF && c == 0 {
				break
			} else {
				return -1, err
			}
		}

		for _, b := range buf[:c] {
			if b == '\n' {
				n++
			}
		}
	}

	if err == io.EOF {
		err = nil
	}
	return n, err
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

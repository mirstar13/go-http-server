package main

import (
	"strings"
)

func cleanInput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "")
	s = strings.TrimSpace(s)
	return s
}

func cleanBuffer(data []byte) []byte {
	for i := 0; i < len(data)-2; i++ {
		if data[i] != 0 && data[i+1] == 0 && data[i+2] == 0 {
			data = data[:i+1]
		}
	}
	return data
}

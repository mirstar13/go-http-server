package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func CleanBuffer(data []byte) []byte {
	for i := 0; i < len(data)-2; i++ {
		if data[i] != 0 && data[i+1] == 0 && data[i+2] == 0 {
			data = data[:i+1]
		}
	}
	return data
}

func GzipEncode(data string) (string, error) {
	var buf bytes.Buffer

	zw := gzip.NewWriter(&buf)

	fmt.Println("data: " + data)

	_, err := zw.Write([]byte(data))
	if err != nil {
		return "", err
	}

	zw.Close()

	return buf.String(), nil
}

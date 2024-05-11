package methods

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mirstar13/go-http-server/utils"
)

const (
	HttpVersion = "HTTP/1.1"
)

const (
	Get  = "GET"
	Post = "POST"
)

const (
	Ok       = "OK"
	NotFound = "Not Found"
	Created  = "Created"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type Method struct {
	Name     string
	Args     string
	Headers  map[string]string
	Encoding string
	Body     string
}

func NewMethod(r *bufio.Reader) (*Method, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	spltLine := strings.Split(line, " ")

	method := &Method{}
	switch spltLine[0] {
	default:

		return nil, ErrMethodNotAllowed

	case Get:
		method.Name = Get
		method.Args = spltLine[1]

		headers, err := getReqHeaders(r)
		if err != nil {
			return nil, err
		}

		method.Headers = headers

		strSize, exists := method.Headers[ContentLength]
		if exists {
			size, err := strconv.Atoi(strSize)
			if err != nil {
				return nil, err
			}

			body, err := getReqBody(r, size)
			if err != nil {
				return nil, err
			}

			method.Body = body
		}
	case Post:
		method.Name = Post
		method.Args = spltLine[1]

		headers, err := getReqHeaders(r)
		if err != nil {
			return nil, err
		}

		method.Headers = headers

		strSize, exists := method.Headers[ContentLength]
		if exists {
			size, err := strconv.Atoi(strSize)
			if err != nil {
				return nil, err
			}

			body, err := getReqBody(r, size)
			if err != nil {
				return nil, err
			}

			method.Body = body
		}
	}

	return method, nil
}

func getReqHeaders(r *bufio.Reader) (map[string]string, error) {
	res := make(map[string]string)

	line := ""
	for {
		var err error

		line, err = r.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("could not read from connection: %w", err)
		}

		if line == "\r\n" {
			break
		}

		line = strings.ReplaceAll(line, "\r\n", "")

		spltLine := strings.Split(line, " ")

		if spltLine[0] == AcceptEncoding+":" {
			encodings := strings.Split(strings.Join(spltLine[1:], " "), ", ")

			for _, encoding := range encodings {
				enc, exist := AcceptedEncodings[encoding]

				if exist {
					if _, ok := res[ContentEncoding]; !ok {
						res[ContentEncoding] += enc
						continue
					}

					res[ContentEncoding] += ", " + enc
				}
			}

			continue
		}

		res[strings.ReplaceAll(spltLine[0], ":", "")] = spltLine[1]
	}

	return res, nil
}

func getReqBody(r *bufio.Reader, size int) (string, error) {
	buf := make([]byte, size)

	_, err := r.Read(buf)
	if err != nil {
		return "", fmt.Errorf("could not read from connection: %w", err)
	}

	res := utils.CleanBuffer(buf)

	return string(res), nil
}

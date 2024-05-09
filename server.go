package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	HTTPVersion = "HTTP/1.1"
)

const (
	MethodNotAllowed = HTTPVersion + " 500 METHOD NOT ALLOWED\r\n\r\n"
	Empty200Resp     = HTTPVersion + " 200 OK\r\n\r\n"
	Empty404Resp     = HTTPVersion + " 404 NOT FOUND\r\n\r\n"
)

type Config struct {
	Flags map[string]*string
}

func main() {
	dirFlag := flag.String("directory", "/", "--directory <directory>")
	flag.Parse()

	cfg := Config{
		Flags: map[string]*string{"dir": dirFlag},
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Println("Could not bind port: " + err.Error())
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: " + err.Error())
			continue
		}

		go cfg.handleRequest(conn)
	}
}

func (cfg *Config) handleRequest(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		err := cfg.parseRequest(r, w)
		if err != nil {
			log.Println("Error parsing request: " + err.Error())
			return
		}
	}
}

func (cfg *Config) parseRequest(r *bufio.Reader, w *bufio.Writer) error {
	line, err := r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("could not read from connection: %w", err)
	}

	spltLine := strings.Split(line, " ")

	resp := ""
	switch spltLine[0] {
	case "GET":
		switch {
		case strings.Contains(spltLine[1], "/files/"):
			fileName := strings.ReplaceAll(spltLine[1], "/files/", "")
			dirToFile := fmt.Sprintf("%s%s", *cfg.Flags["dir"], fileName)

			if _, err := os.Stat(dirToFile); err == nil {
				file, err := os.Open(dirToFile)
				if err != nil {
					return fmt.Errorf("could not open specified file: %w", err)
				}
				defer file.Close()

				data, err := readFile(file)
				if err != nil {
					return fmt.Errorf("could not read from file: %w", err)
				}

				resp = formatResp(string(data), 200, "OK", "application/octet-stream")

			} else {
				resp = formatResp("", 404, "NOT FOUND", "application/octet-stream")
			}

		case strings.Contains(spltLine[1], "echo"):
			data := strings.ReplaceAll(spltLine[1], "/echo/", "")
			resp = formatResp(data, 200, "OK", "text/plain")

		case spltLine[1] == "/user-agent":
			r.ReadString('\n') // Skip Host header

			usrAgentHeader, err := r.ReadString('\n')
			if err != nil {
				return err
			}

			spltHeader := strings.Split(usrAgentHeader, " ")
			usrAgent := cleanInput(spltHeader[1])
			resp = formatResp(usrAgent, 200, "OK", "text/plain")

		default:
			if spltLine[1] == "/" {
				resp = Empty200Resp
			} else {
				resp = Empty404Resp
			}
		}
	case "POST":
		switch {
		case strings.Contains(spltLine[1], "/files/"):
			filename := strings.ReplaceAll(spltLine[1], "/files/", "")
			dirToFile := fmt.Sprintf("%s%s", *cfg.Flags["dir"], filename)

			fmt.Println(dirToFile)

			file, err := os.Create(dirToFile)
			if err != nil {
				return fmt.Errorf("could not create file: %w", err)
			}

			fmt.Println("created file")

			for { // Skip headers
				line, err = r.ReadString('\n')
				if err != nil {
					return fmt.Errorf("could not read from connection: %w", err)
				}
				fmt.Println(line)
				if cleanInput(line) == "" {
					break
				}
			}

			fmt.Println("skiped headers")

			data := make([]byte, 256)
			r.Read(data)

			data = cleanBuffer(data)
			data = []byte(cleanInput(string(data)))

			err = writeToFile(file, data)
			if err != nil {
				return err
			}

			resp = formatResp("", 201, "", "application/octet-stream")

		}
	}

	if resp != "" {
		_, err = w.Write([]byte(resp))
		if err != nil {
			return fmt.Errorf("cound not write to connection: %w", err)
		}
		w.Flush()
	}

	return nil
}

func formatResp(data string, status int, msg string, content string) string {
	return fmt.Sprintf("%s %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", HTTPVersion, status, msg, content, len(data), data)
}

func writeToFile(file *os.File, data []byte) error {
	_, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}

func readFile(file *os.File) ([]byte, error) {
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
	}
	return cleanBuffer(buf), nil
}

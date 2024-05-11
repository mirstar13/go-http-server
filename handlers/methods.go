package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mirstar13/go-http-server/methods"
	"github.com/mirstar13/go-http-server/utils"
)

func handleGet(h *Handler, method *methods.Method) error {
	headers := []string{}

	if encodings, ok := method.Headers[methods.ContentEncoding]; ok {
		headers = append(headers, methods.FormatHeader(methods.ContentEncoding, encodings))
	}

	switch {
	default:
		resp := formatEmptyResp(404, methods.NotFound)

		h.writer.WriteString(resp)

	case method.Args == "/":
		resp := formatEmptyResp(200, methods.Ok)

		h.writer.WriteString(resp)

	case strings.Contains(method.Args, "/echo/"):
		data := strings.ReplaceAll(method.Args, "/echo/", "")

		if enc, ok := method.Headers[methods.ContentEncoding]; ok {
			switch enc {
			case methods.Gzip:
				var err error

				data, err = utils.GzipEncode(data)
				if err != nil {
					return fmt.Errorf("could not encode data: %w", err)
				}

			}
		}

		headers = append(headers, methods.FormatContentTypeHeader(methods.TextPlain, len(data))...)
		resp := formatRespWithContent(200, methods.Ok, data, headers...)

		h.writer.WriteString(resp)

	case method.Args == "/user-agent":
		data := method.Headers[methods.UserAgent]

		headers = append(headers, methods.FormatContentTypeHeader(methods.TextPlain, len(data))...)
		resp := formatRespWithContent(200, methods.Ok, data, headers...)

		h.writer.WriteString(resp)

	case strings.Contains(method.Args, "/files/"):
		fileName := strings.ReplaceAll(method.Args, "/files/", "")
		pathToFile := h.cfg.Dir() + fileName

		if _, err := os.Stat(pathToFile); err == nil {
			file, err := os.Open(pathToFile)
			if err != nil {
				return fmt.Errorf("could not open file: %w", err)
			}
			defer file.Close()

			strSize := method.Headers[methods.ContentLength]
			if strSize == "" {
				strSize = "1024"
			}

			size, err := strconv.Atoi(strSize)
			if err != nil {
				return err
			}

			data, err := readFile(file, size)
			if err != nil {
				return err
			}

			headers = append(headers, methods.FormatContentTypeHeader(methods.AppOctetStream, len(data))...)
			resp := formatRespWithContent(200, methods.Ok, data, headers...)

			h.writer.WriteString(resp)

		} else {
			headers = append(headers, methods.FormatContentTypeHeader(methods.AppOctetStream, 0)...)
			resp := formatRespWithContent(404, methods.NotFound, "", headers...)

			h.writer.WriteString(resp)

		}

	}

	return nil
}

func handlePost(h *Handler, method *methods.Method) error {
	switch {
	default:
		resp := formatEmptyResp(404, methods.NotFound)

		h.writer.WriteString(resp)

	case strings.Contains(method.Args, "/files/"):
		fileName := strings.ReplaceAll(method.Args, "/files/", "")
		pathToFile := h.cfg.Dir() + fileName

		file, err := os.Create(pathToFile)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		defer file.Close()

		err = writeFile(file, method.Body)
		if err != nil {
			return err
		}

		resp := formatEmptyResp(201, methods.Created)
		h.WriteResponse(resp)

	}

	return nil
}

func formatRespWithContent(status int, msg string, body string, headers ...string) string {
	resp := fmt.Sprintf("%s %d %s\r\n", methods.HttpVersion, status, msg)

	for _, header := range headers {
		resp += header + "\r\n"
	}
	resp += "\r\n"

	resp += body

	return resp
}

func formatEmptyResp(status int, msg string) string {
	return fmt.Sprintf("%s %d %s\r\n\r\n", methods.HttpVersion, status, msg)
}

func readFile(file *os.File, size int) (string, error) {
	buf := make([]byte, size)

	_, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("could not read from file: %w", err)
	}

	buf = utils.CleanBuffer(buf)

	return string(buf), nil
}

func writeFile(file *os.File, data string) error {
	_, err := file.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

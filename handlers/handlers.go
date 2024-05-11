package handlers

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/mirstar13/go-http-server/methods"
	"github.com/mirstar13/go-http-server/server/config"
)

type Handler struct {
	connection net.Conn
	cfg        *config.Config
	reader     *bufio.Reader
	writer     *bufio.Writer
}

var methodHandlers = map[string]func(*Handler, *methods.Method) error{
	methods.Get:  handleGet,
	methods.Post: handlePost,
}

func NewHandler(conn net.Conn, cfg *config.Config) *Handler {
	return &Handler{
		connection: conn,
		cfg:        cfg,
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
	}
}

func (h *Handler) HandlerClient() error {
	defer h.connection.Close()

	for {
		req, err := methods.NewMethod(h.reader)
		if err != nil {
			return fmt.Errorf("failed to parse req: %v", err)
		}

		err = h.handleMethod(req)
		if err != nil {
			return err
		}

		h.writer.Flush()
	}
}

func (h *Handler) WriteResponse(msg string) {
	h.writer.WriteString(msg)
}

func (h *Handler) handleMethod(method *methods.Method) error {
	intstraction := strings.TrimSpace(method.Name)
	handler, exists := methodHandlers[intstraction]

	if !exists {
		return fmt.Errorf("unknown method: %s", intstraction)
	}

	return handler(h, method)
}

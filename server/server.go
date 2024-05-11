package server

import (
	"fmt"
	"log"
	"net"

	"github.com/mirstar13/go-http-server/handlers"
	"github.com/mirstar13/go-http-server/server/config"
)

type Server struct {
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (srv *Server) Start() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+srv.config.Port())
	if err != nil {
		return fmt.Errorf("could not bind port "+srv.config.Port()+": %w", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: " + err.Error())
			continue
		}

		connHandler := handlers.NewHandler(conn, srv.config)

		go srv.serveConnection(connHandler)
	}
}

func (s *Server) serveConnection(connHandler *handlers.Handler) {
	err := connHandler.HandlerClient()
	if err != nil {
		log.Printf("something happend with the client connection: " + err.Error())
	}
}

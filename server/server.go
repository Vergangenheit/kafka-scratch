package server

import (
	"fmt"
	"net"

	"github.com/hashicorp/go-hclog"
)

const (
	receiveBuf = 1024
)

type Server struct {
	host   string
	port   string
	Logger hclog.Logger
}

func NewServer(host string, port string, logger hclog.Logger) *Server {
	return &Server{
		host:   host,
		port:   port,
		Logger: logger,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	s.Logger.Info("Connected to client:", conn.RemoteAddr())

	data, err := s.parseRequest(conn)
	if err != nil {
		s.Logger.Error("Error parsing request:", err)
		return
	}
	s.Logger.Info("Received request: ", data)

}

func (s *Server) generateResponse(data string) ([][]byte, error) {
	resp := make([][]byte, 0)
	// get correlation id from data
	// it starts after 16 characters
	if len(data) < 24 {
		return nil, fmt.Errorf("data too short")
	}
	correlationID := data[16:24]
	// send it back to client
	corrIdBytes := []byte(correlationID)
	resp = append(resp, corrIdBytes)
	return resp, nil
}

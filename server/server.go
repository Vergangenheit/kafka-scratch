package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	var buffer bytes.Buffer
	reader := bufio.NewReader(conn)
	for {
		chunk := make([]byte, 1024) // Temporary buffer for reading chunks
		n, err := reader.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			s.Logger.Error("Error parsing response:", err)
			return
		}
		_, err = buffer.Write(chunk[:n])
		if err != nil {
			if err == bytes.ErrTooLarge {
				break
			}
			s.Logger.Error("Error writing to buffer:", err)
			return
		}

	}
	// Process the complete received data
	data := buffer.Bytes()
	dataStr := string(data)
	s.Logger.Info("Received request:", dataStr)

}

func (s *Server) generateResponse(conn net.Conn) ([][]byte, error) {
	// buffer the conn
	buffer := make([]byte, receiveBuf)
	// start reading chunks delimited by newline byte

	_, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	// parse the request
	parsedString := parseBulkBytes(buffer)
	s.Logger.Info("Received request:", parsedString)
	return [][]byte{
		{0, 0, 0, 0, 0, 0, 0, 7},
	}, nil
}

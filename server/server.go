package server

import (
	"fmt"
	"io"
	"net"
	"os"

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
	for {
		responses, err := s.generateResponse(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			s.Logger.Error("Error parsing response:", err)
			return
		}
		for _, resp := range responses {
			_, err = conn.Write(resp)
			if err != nil {
				s.Logger.Error("Error sending response:", err)
				return
			}
		}
	}

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	srv := NewServer("localhost", "9092", hclog.Default())
	err := srv.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

func (s *Server) generateResponse(conn net.Conn) ([][]byte, error) {
	// buffer the conn
	buffer := make([]byte, receiveBuf)
	// start reading chunks delimited by newline byte

	_, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return [][]byte{
		{0, 0, 0, 0, 0, 0, 0, 7},
	}, nil
}

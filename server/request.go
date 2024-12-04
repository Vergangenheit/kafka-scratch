package server

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strings"
)

func parseBulkBytes(input []byte) []string {
	// Convert byte slice to string
	str := string(input)

	// Trim any trailing \r\n to prevent an extra empty element after splitting
	str = strings.TrimSuffix(str, "\r\n")

	// Split the string by "\r\n"
	parts := strings.Split(str, "\r\n")

	return parts
}

func (s *Server) parseRequest(conn net.Conn) (string, error) {
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
			return "", err
		}
		_, err = buffer.Write(chunk[:n])
		if err != nil {
			if err == bytes.ErrTooLarge {
				break
			}
			s.Logger.Error("Error writing to buffer:", err)
			return "", err
		}

	}
	// Process the complete received data
	data := buffer.Bytes()
	dataStr := string(data)

	return dataStr, nil
}

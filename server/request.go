package server

import "strings"

func parseBulkBytes(input []byte) []string {
	// Convert byte slice to string
	str := string(input)

	// Trim any trailing \r\n to prevent an extra empty element after splitting
	str = strings.TrimSuffix(str, "\r\n")

	// Split the string by "\r\n"
	parts := strings.Split(str, "\r\n")

	return parts
}

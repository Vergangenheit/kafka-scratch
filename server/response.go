package server

type Header struct {
	CorrelationID int32
}

type Response struct {
	MessageSize int32
	Header      Header
}

package client

import "net"

type KafkaClient struct {
	conn net.Conn
}

func NewKafkaClient(serverAddress string) (*KafkaClient, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &KafkaClient{conn: conn}, nil
}

func (r *KafkaClient) Send(request []byte) error {
	_, err := r.conn.Write([]byte(request))
	if err != nil {
		return err
	}
	return nil
}

func (r *KafkaClient) Close() {
	r.conn.Close()
}

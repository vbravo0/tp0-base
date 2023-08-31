package common

import (
	"encoding/binary"
	"errors"
	"net"
)

const SIZE_U32 = 4

func sendBytes(socket net.Conn, msg []byte) error {
	totalSent := 0
	for totalSent < len(msg) {
		sent, err := socket.Write(msg[totalSent:])
		if err != nil {
			return errors.New("socket send error")
		}
		totalSent += sent
	}
	return nil
}

func recvBytes(socket net.Conn, size int) ([]byte, error) {
	buffer := make([]byte, 1024)
	dataReceivedCounter := 0
	for dataReceivedCounter < size {
		n, err := socket.Read(buffer[dataReceivedCounter:])
		if err != nil {
			return buffer, errors.New("socket recv error")
		}
		dataReceivedCounter += n
	}
	return buffer, nil
}

func sendU32(socket net.Conn, n uint32) error {
	bs := make([]byte, SIZE_U32)
	binary.BigEndian.PutUint32(bs, n)
	return sendBytes(socket, bs)
}

func recvU32(socket net.Conn) (uint32, error) {
	data, err := recvBytes(socket, SIZE_U32)
	if err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint32(data)
	return n, nil
}

package common

import (
	"encoding/binary"
	"errors"
	"fmt"
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

func recvBytes(socket net.Conn, size uint32) ([]byte, error) {
	buffer := make([]byte, size)
	var dataReceivedCounter uint32 = 0
	for dataReceivedCounter < size {
		n, err := socket.Read(buffer[dataReceivedCounter:])
		if err != nil {
			return buffer, errors.New("socket recv error")
		}
		dataReceivedCounter += uint32(n)
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

func sendString(socket net.Conn, s string) error {
	size := uint32(len(s))
	err := sendU32(socket, size)
	if err != nil {
		return err
	}

	data := []byte(s)
	err = sendBytes(socket, data)
	if err != nil {
		return err
	}

	return nil
}

func recvString(socket net.Conn) (string, error) {
	size, err := recvU32(socket)
	if err != nil {
		return "", err
	}

	data, err := recvBytes(socket, size)
	if err != nil {
		return "", err
	}

	return string(data[:]), nil
}

func sendBet(socket net.Conn, name string, surname string, document string, birthdate string, number int) error {
	res := fmt.Sprintf("%v,%v,%v,%v,%v", name, surname, document, birthdate, number)
	return sendString(socket, res)
}

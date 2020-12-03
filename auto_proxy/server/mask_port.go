package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

func EncodePort(port string) (encode string, err error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return
	}

	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(portInt))

	for i, b := range bytes {
		encode += ByteTo16(b ^ portMask[i])
	}

	return
}

func ByteTo16(b byte) string {
	return fmt.Sprintf("%02x", b)
}

func DecodePortByBytes(portBytes []byte) (encode string, err error) {
	if len(portBytes) != 2 {
		return "", errors.New("invalid port bytes")
	}

	decodeBytes := make([]byte, 2)
	for i, b := range portBytes {
		decodeBytes[i] += b ^ portMask[i]
	}

	return strconv.Itoa(int(binary.LittleEndian.Uint16(decodeBytes))), nil
}

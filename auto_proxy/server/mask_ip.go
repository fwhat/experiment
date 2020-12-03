package server

import (
	"encoding/hex"
	"errors"
	"net"
)

// Is p all zeros?
func isZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}

func isIpV4(ip net.IP) bool {
	if len(ip) == net.IPv4len {
		return true
	}

	if len(ip) == net.IPv6len && isZeros(ip[0:10]) {
		return true
	}

	return false
}

func encodeIpv4(ip net.IP) (encode string) {
	for i, b := range ip.To4() {
		encode += ByteTo16(b ^ ip4Mask[i])
	}

	return
}

func EncodeIp(ip net.IP) (encode string, err error) {
	if isIpV4(ip) {
		return encodeIpv4(ip), nil
	}

	return encode, errors.New("invalid ip")
}

func encodeIpWithLen(ip net.IP) (encode string, err error) {
	encode, err = EncodeIp(ip)
	if err != nil {
		return "", err
	}

	return ByteTo16(byte(len(encode)/2)^domainLenMask) + encode, nil
}

func DecodeIp(encode string) (ip net.IP, err error) {
	b, err := hex.DecodeString(encode)
	if err != nil {
		return
	}

	return decodeIpByBytes(b)
}

func decodeIpByBytes(b []byte) (ip net.IP, err error) {
	if len(b) == net.IPv4len {
		ip = make([]byte, 4)
		for i, v := range b {
			ip[i] = v ^ ip4Mask[i]
		}

		return ip, nil
	}

	return ip, errors.New("invalid encode str")
}

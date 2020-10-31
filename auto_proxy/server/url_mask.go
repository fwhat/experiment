package server

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

var schemeMask = []byte{0x33, 0x12, 0x45, 0x89, 0x11, 0x44}
var ip4Mask = []byte{0x60, 0x37, 0x58, 0x4c}
var portMask = []byte{0x21, 0x68}

var schemeLenMask = byte(0x48)
var domainLenMask = byte(0x99)

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

func DecodeUrl(encode string) (u *url.URL, err error) {
	b, err := hex.DecodeString(encode)

	if err != nil {
		return
	}
	u = &url.URL{}

	schemeLen := b[0] ^ schemeLenMask

	if int(schemeLen) > len(schemeMask) {
		return nil, errors.New("invalid encode str")
	}

	for i, v := range b[1 : schemeLen+1] {
		u.Scheme += string(schemeMask[i] ^ v)
	}

	domainLen := b[schemeLen+1] ^ domainLenMask
	domain, err := decodeIpByBytes(b[schemeLen+2 : domainLen+schemeLen+2])
	if err != nil {
		return nil, err
	}

	portStr, err := DecodePortByBytes(b[domainLen+schemeLen+2:])
	if err != nil {
		return
	}

	u.Host = domain.String() + ":" + portStr

	return
}

func EncodeUrl(u *url.URL) (encode string, err error) {
	encode, err = encodeSchemeWithLen(u.Scheme)
	if err != nil {
		return
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", err
	}

	var domainEncode = ""

	parseIP := net.ParseIP(host)
	if parseIP == nil {
		// domain
		return "", errors.New("domain is not support")
	} else {
		domainEncode, err = encodeIpWithLen(parseIP)
		if err != nil {
			return "", err
		}

		encode += domainEncode
	}

	encodePort, err := EncodePort(port)
	if err != nil {
		return "", err
	}

	return encode + encodePort, nil
}

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

func encodeScheme(scheme string) (encode string, err error) {
	if len(scheme) > len(schemeMask) {
		return "", errors.New("not support scheme")
	}

	for i, b := range scheme {
		encode += ByteTo16(byte(b) ^ schemeMask[i])
	}

	return
}

func encodeSchemeWithLen(scheme string) (encode string, err error) {
	encode, err = encodeScheme(scheme)
	if err != nil {
		return
	}

	return ByteTo16(byte(len(encode)/2)^schemeLenMask) + encode, nil
}

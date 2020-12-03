package server

import (
	"encoding/hex"
	"errors"
	"net"
	"net/url"
)

var schemeMask = []byte{0x33, 0x12, 0x45, 0x89, 0x11, 0x44}
var ip4Mask = []byte{0x60, 0x37, 0x58, 0x4c}
var portMask = []byte{0x21, 0x68}
var domainMask = []byte{0x19, 0x23, 0x99, 0x33, 0x41, 0x91, 0x10, 0x24, 0x56, 0x78, 0x98, 0x88, 0x45, 0x54, 0x18, 0x71, 0x93, 0x21, 0x45, 0x78, 0x81, 0x51, 0x81}

var schemeLenMask = byte(0x48)
var domainLenMask = byte(0x99)

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
		domainEncode, err = encodeDomainWithLen(host)
		if err != nil {
			return "", err
		}
	} else {
		domainEncode, err = encodeIpWithLen(parseIP)
		if err != nil {
			return "", err
		}
	}

	encode += domainEncode

	encodePort, err := EncodePort(port)
	if err != nil {
		return "", err
	}

	return encode + encodePort, nil
}

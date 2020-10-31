package server

import (
	"net"
	"net/url"
	"testing"
)

func TestEncodeIp4(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")

	hash, err := EncodeIp(ip)
	if err != nil {
		t.Error(err)
	}

	decodeIp, err := DecodeIp(hash)
	if err != nil {
		t.Error(err)
	}

	if !decodeIp.Equal(ip) {
		t.Error("decode value not equal origin")
	}
}

func TestEncodeUrl(t *testing.T) {
	parse, err := url.Parse("tcp://127.0.0.1:3366")
	if err != nil {
		t.Error(err)
	}

	encodeUrl, err := EncodeUrl(parse)
	if err != nil {
		t.Error(err)
	}

	decodeUrl, err := DecodeUrl(encodeUrl)
	if err != nil {
		t.Error(err)
	}
	if parse.String() != decodeUrl.String() {
		t.Error("decode url fail")
	}
}

func TestByteTo16(t *testing.T) {
	if ByteTo16(byte(7)) != "07" {
		t.Error("len error")
	}
}

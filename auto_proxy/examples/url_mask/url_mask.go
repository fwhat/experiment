package main

import (
	"auto_proxy/server"
	"net/url"
)

func main() {
	parse, err := url.Parse("tcp://127.0.0.1:5500")
	if err != nil {
		panic(err)
	}

	encodeUrl, err := server.EncodeUrl(parse)

	if err != nil {
		panic(err)
	}
	println(encodeUrl)
}

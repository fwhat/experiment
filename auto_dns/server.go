package main

import (
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

func StartServer() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)

		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Println(err)
			continue
		}
		go ServerDNS(addr, conn, msg)
	}
}

func ServerDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	// query info
	if len(msg.Questions) < 1 {
		return
	}

	question := msg.Questions[0]

	fmt.Printf("[%s] queryName: [%s]\n", question.Type.String(), question.Name.String())

	packed, err := msg.Pack()
	if err != nil {
		panic(err)
	}

	dial, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	_, err = dial.Write(packed)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 512)
	_, err = dial.Read(buf)
	if err != nil {
		panic(err)
	}
	var msg2 dnsmessage.Message
	if err := msg2.Unpack(buf); err != nil {
		panic(err)
	}
	if msg2.Header.Response {
		packed, err = msg2.Pack()
		_, err = conn.WriteToUDP(packed, addr)
	}
}

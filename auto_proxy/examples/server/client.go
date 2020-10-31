package main

import (
	"auto_proxy/server"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	go func() {
		panic(proxyServer())
	}()

	go func() {
		panic(realServer())
	}()

	// tcp://127.0.0.1:5500
	dial, err := net.Dial("tcp", "4b4771359d1f37584d5d7d.localhost:4400")
	if err != nil {
		panic(err)
	}
	defer dial.Close()

	go func() {
		for {
			buf := make([]byte, 11)
			_, err := dial.Read(buf)
			if err != nil {
				panic(err)
			}

			log.Println("recv: " + string(buf))
		}
	}()

	for i := 0; i < 100; i++ {
		dial.Write([]byte("hello world"))
	}

	<-make(chan interface{})
}

func proxyServer() error {
	return server.NewServer().Serve(":4400")
}

func realServer() error {
	listen, err := net.Listen("tcp", ":5500")
	if err != nil {
		panic(err)
	}

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := listen.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Println("accept tcp error")
				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		go func() {
			_, err := io.Copy(conn, conn)
			if err != nil {
				panic(err)
			}
		}()
	}
}

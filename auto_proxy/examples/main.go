package main

import (
	"io"
	"log"
	"net"
	"time"
)

func server() error {
	listen, err := net.Listen("tcp", ":3322")
	if err != nil {
		return err
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
				log.Println(err)
				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		go func() {
			defer conn.Close()

			distConn, err := net.DialTimeout("tcp", "www.baidu.com:80", time.Second*3)
			if err != nil {
				panic(err)
			}

			defer distConn.Close()

			go func() {
				_, err = io.Copy(distConn, conn)
				_ = distConn.Close()
				log.Printf("server handle conn to distConn error %v\n", err)
				return
			}()

			_, err = io.Copy(conn, distConn)
			if err != nil {
				log.Printf("server handle conn to srcConn error %v\n", err)
			}
		}()
	}
}

func main() {
	panic(server())
}

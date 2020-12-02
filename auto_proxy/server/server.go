package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve(addr string) error {
	listen, err := net.Listen("tcp", addr)
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
			err := s.handleConn(conn)
			if err != nil && err != io.EOF {
				log.Println(err)
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	servername := getSNIServerName(buf[:n])
	if servername == "" {
		return conn.Close()
	}

	return s.proxy(conn, buf[:n], servername)
}

func getSNIServerName(buf []byte) string {
	n := len(buf)
	if n <= 5 {
		log.Println("not tls handshake")
		return ""
	}

	// tls record type
	if recordType(buf[0]) != recordTypeHandshake {
		log.Println("not tls")
		return ""
	}

	// tls major version
	if buf[1] != 3 {
		log.Println("TLS version < 3 not supported")
		return ""
	}

	// payload length
	//l := int(buf[3])<<16 + int(buf[4])

	//log.Printf("length: %d, got: %d", l, n)

	// handshake message type
	if buf[5] != typeClientHello {
		log.Println("not client hello")
		return ""
	}

	// parse client hello message

	msg := &clientHelloMsg{}

	// client hello message not include tls header, 5 bytes
	ret := msg.unmarshal(buf[5:n])
	if !ret {
		log.Println("parse hello message return false")
		return ""
	}
	return msg.serverName
}

func (s *Server) proxy(srcConn net.Conn, data []byte, distName string) error {
	defer srcConn.Close()

	split := strings.Split(distName, ".")
	if len(split) < 2 {
		return errors.New("not allow")
	}

	url, err := DecodeUrl(split[0])
	if err != nil {
		return errors.New(fmt.Sprintf("not allow %v", err))
	}

	distConn, err := net.DialTimeout(url.Scheme, url.Host, time.Second*3)

	go func() {
		_, err = io.Copy(distConn, srcConn)
		_ = distConn.Close()
		log.Printf("server handle conn to distConn error %v\n", err)
		return
	}()

	_, err = distConn.Write(data)
	if err != nil {
		return err
	}

	_, err = io.Copy(srcConn, distConn)
	if err != nil {
		log.Printf("server handle conn to srcConn error %v\n", err)
	}

	return nil
}

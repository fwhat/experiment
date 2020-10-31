package server

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		http.NotFound(w, r)
	}

	split := strings.Split(host, ".")
	if len(split) < 2 {
		http.NotFound(w, r)
		return
	}

	url, err := DecodeUrl(split[0])
	if err != nil {
		Error(w, "url is invalid")
		return
	}

	distConn, err := net.DialTimeout(url.Scheme, url.Host, time.Second*3)
	if err != nil {
		Error(w, "dial remote timeout")
		return
	}

	srcConn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		Error(w, "hijack fail")
		return
	}

	defer srcConn.Close()

	go func() {
		_, err = io.Copy(distConn, srcConn)
		_ = distConn.Close()
		log.Printf("server handle conn to distConn error %v\n", err)
		return
	}()

	_, err = io.Copy(srcConn, distConn)

	log.Printf("server handle conn to srcConn error %v\n", err)
}

func Error(w http.ResponseWriter, error string) {
	w.WriteHeader(500)
	w.Write([]byte(error))
}

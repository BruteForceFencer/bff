// Package server implements a server following its own custom protocol.  The
// protocol works as follows.
//
//     1. A client connects and sends an object in JSON format matching the
//     structure of the Request type.
//
//     2. The server returns a single character. "t" means that the request is
//     valid.  "f" means that the request is invalid (either an attack or an
//     error).
//
//     3. The client disconnects or goes again from step 1.
package controlserver

import (
	"log"
	"net"
	"os"
	"strings"
)

// Server is a server that interprets requests according to the protocol.
type Server struct {
	HandleFunc      func(*Request) bool
	AcceptedSources map[string]bool
	listener        net.Listener
}

func New() *Server {
	return &Server{
		AcceptedSources: map[string]bool{
			"":          true,
			"127.0.0.1": true,
		},
	}
}

// Blocks and listens for requests.
func (s *Server) ListenAndServe(typ, addr string) error {
	// Remove any old socket.
	if typ == "unix" {
		os.Remove(addr)
	}

	// Start listening.
	var err error
	s.listener, err = net.Listen(typ, addr)
	if err != nil {
		return err
	}

	// Accept requests.
	go s.acceptRequests()
	return nil
}

func (s *Server) acceptRequests() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("connection error:", err)
			continue
		}

		if !s.accepted(conn) {
			conn.Close()
			continue
		}

		// For performance, we launch every handler in its own goroutine.
		go func(conn net.Conn) {
			for {
				request, err := ReadRequest(conn)
				if err != nil {
					conn.Close()
					return
				}

				response := s.HandleFunc(request)
				if response {
					conn.Write([]byte("t"))
				} else {
					conn.Write([]byte("f"))
				}
			}
		}(conn)
	}
}

// accepted returns true if the connection comes from a trusted source.
func (s *Server) accepted(conn net.Conn) bool {
	// Separate the ip address from the port.
	addr := conn.RemoteAddr().String()
	colonIndex := strings.Index(addr, ":")

	// If there's no colon, then this request comes from a Unix socket.
	if colonIndex == -1 {
		allowed, ok := s.AcceptedSources[addr]
		return allowed && ok
	}

	// Remove the port number from the end of the address.
	addr = addr[0:colonIndex]
	allowed, ok := s.AcceptedSources[addr]
	return allowed && ok
}

// Close stops the server.
func (s *Server) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	s.listener.Close()
}

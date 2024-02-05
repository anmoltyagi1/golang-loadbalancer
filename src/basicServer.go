package main

import (
	"net/http/httputil"
	"net/http"
)

// define a server struct here
type basicServer struct {
	address string
	connections int
	proxy   *httputil.ReverseProxy
}


// address method
func (s *basicServer) Address() string {
	return s.address 
}

// isAlive method - returning true for sake of simplicity 
func (s *basicServer) isAlive() bool {
	return true
}

func (s *basicServer) Connections() int {
	// fmt.Printf("Connections: %s\n", s.address)
	return s.connections
}

func (s *basicServer) SetConnections(newConnection int) int {
    s.connections = newConnection
    return s.connections
}


// Serve method
func (s *basicServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)	
}
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)


// define a server interface here
type Server interface {
	Address() string
	isAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}


// define a server struct here
type basicServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

// define the load balancer struct here
type LoadBalancer struct {
	port 		  string
	roundRobinIndex int
	servers         []Server

}

// initialize the load balancer here
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		roundRobinIndex: 0,
		servers: servers,
		port: port,
	}
}

// implement the new server method here
func newServer(address string) *basicServer {
	serverUrl, err := url.Parse(address) 
 
	handleErr(err)

	return &basicServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),

	}
}


// handle any errors
func handleErr(err error){
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}


// address method
func (s *basicServer) Address() string {
	return s.address 
}

// isAlive method - returning true for sake of simplicity 
func (s *basicServer) isAlive() bool {
	return true
}

// Serve method
func (s *basicServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// getNextAvailableServer method
func (lb *LoadBalancer) getNextAvailableServer() Server {
    numServers := len(lb.servers)


	// round robin being applied - go to the next one if current not available
    for i := 0; i < numServers; i++ {
        server := lb.servers[lb.roundRobinIndex]
        lb.roundRobinIndex = (lb.roundRobinIndex + 1) % numServers

        if server.isAlive() {
            return server
        }
    }

    return nil
}

// serveProxy method
func(lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	target := lb.getNextAvailableServer()
	fmt.Printf("Redirecting request to: %s\n", target.Address())
	target.Serve(w, r)
}

// main function
func main(){

	// get all the servers
	servers := []Server{
		newServer("https://www.meta.com"),
		newServer("https://www.google.com"),
		newServer("https://www.stripe.com"),
	}
	

	// create a new load balancer
	lb := NewLoadBalancer("8000", servers)

	// do the redirection 
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}


	// start the load balancer
	http.HandleFunc("/", handleRedirect)


	fmt.Printf("Load Balancer started at :%s\n", lb.port)

	// listen
	http.ListenAndServe(":"+lb.port, nil)
}
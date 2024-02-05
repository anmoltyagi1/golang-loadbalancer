package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"container/heap"
	"encoding/json"
)


// define a server interface here
type Server interface {
	Address() string
	isAlive() bool
	Connections() int
	SetConnections (int) int
	Serve(rw http.ResponseWriter, req *http.Request)
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}


// initialize the load balancer here
func NewLoadBalancer(port string, servers []Server, minHeap *MinHeap, mySet Set) *LoadBalancer {
	return &LoadBalancer{
		roundRobinIndex: 0,
		servers: servers,
		minHeap: minHeap,
		set: mySet,
		port: port,
	}
}

// implement the new server method here
func newServer(address string) *basicServer {
	serverUrl, err := url.Parse(address) 
 
	handleErr(err)

	return &basicServer{
		address: address,
		connections: 0,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),

	}
}


func checkServerHealth(server Server, lb *LoadBalancer) {
	// Make GET request to the server
	resp, err := http.Get(server.Address())
	if err != nil {
		fmt.Printf("Error checking server health: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		fmt.Printf("Error decoding JSON response: %v\n", err)
		return
	}


	// Check server status
	if healthResp.Status == "DEGRADED" && lb.set.Contains(server.Address()) {
		lb.mutex.Lock()
		defer lb.mutex.Unlock()
		lb.set.Remove(server.Address())
		for i, serverX := range *lb.minHeap {
			if server == serverX[0].(Server) {
				heap.Remove(lb.minHeap, i)
				break
			}
		}
	} else if healthResp.Status == "HEALTHY" && !lb.set.Contains(server.Address()) {
		lb.mutex.Lock()
		defer lb.mutex.Unlock()
		lb.set.Add(server.Address())
		heap.Push(lb.minHeap, []interface{}{server, server.Connections()})
	}
}


// main function
func main(){

	// get all the servers
	servers := []Server{
		newServer("http://localhost:3000/health"),
		newServer("http://localhost:3001/health"),
		newServer("http://localhost:3002/health"),
		newServer("http://localhost:3003/health"),
	}

	minHeap := &MinHeap{}

	for _, server := range servers {
		// Create a slice with address and connections
		serverInfo := []interface{}{server, server.Connections()}
		// Push the slice onto the maxHeap
		heap.Push(minHeap, serverInfo)
	}

	mySet := NewSet()

	// create a set data structure 
	mySet.Add("http://localhost:3000/health")
	mySet.Add("http://localhost:3001/health")
	mySet.Add("http://localhost:3002/health")
	mySet.Add("http://localhost:3003/health")

	// create a new load balancer
	lb := NewLoadBalancer("8000", servers, minHeap, mySet)

	// do the redirection 
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("request : %s", r);
		lb.serveProxy(w, r, minHeap)
	}


	// Run a timer to check server health every 10 seconds
	go func() {
		for {
			// Your existing code for health checks and goroutine launching
			for _, server := range servers {
				go checkServerHealth(server, lb)
			}
	
			fmt.Printf("\n")
	
			// Sleep for 10 seconds
			printSet(lb.set)

			time.Sleep(2 * time.Second)
		}
	}()



	// start the load balancer
	http.HandleFunc("/", handleRedirect)


	fmt.Printf("Load Balancer started at :%s\n", lb.port)

	// listen
	http.ListenAndServe(":"+lb.port, nil)
}




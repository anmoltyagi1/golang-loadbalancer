package main

import (
	"fmt"
	"net/http"
	"container/heap"
	"sync"
)



// define the load balancer struct here
type LoadBalancer struct {
	port 		  string
	roundRobinIndex int
	servers         []Server
	minHeap         *MinHeap
	set             Set
	mutex		   sync.Mutex
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

func (lb *LoadBalancer) routeToLeastConnection() Server {

    // Pop the minimum server from the min-heap
    lb.mutex.Lock()

    serverToConnection := heap.Pop(lb.minHeap)
    serverConnection := serverToConnection.([]interface{})
    server := serverConnection[0].(Server)
    connection := serverConnection[1].(int)
	// fmt.Printf("Server: %s, Connections: %d\n", server.Address(), connection)

	// fmt.Printf("%s ", lb.minHeap);


	heap.Push(lb.minHeap, []interface{}{server, connection + 1})
	fmt.Printf("\n")
	// fmt.Printf("%s ", lb.minHeap);
	// loop through minHeap and print the server address and then the connectoins
	fmt.Printf("STARTTT ----------------\n")

	for _, server := range *lb.minHeap {
		fmt.Printf("Server: %s, Connections: %d\n", server[0].(Server).Address(), server[1].(int))
	}


	lb.mutex.Unlock();



    return server
}


// serveProxy method
func(lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request, minHeap *MinHeap) {

	target := lb.routeToLeastConnection()

	fmt.Printf("\n")

	target.Serve(w, r)

	lb.mutex.Lock();
	

	index := 0
	for i, server := range *lb.minHeap {
		if target == server[0].(Server) {
			// fmt.Printf("found at index: %d\n", i)
			index = i
			break
		}
	}

	(*lb.minHeap)[index][1] = target.Connections() 

	heap.Fix(lb.minHeap, index)

	lb.mutex.Unlock()
	fmt.Printf("\n")

}

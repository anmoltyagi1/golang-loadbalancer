# Load Balancer in Go - Least Connections

This is a basic load balancer implemented in Go that distributes incoming requests to a list of servers using the Least Connections algorithm.

## How It Works

The load balancer listens for incoming requests and routes them to different backend servers. Here's how it functions:

1. The load balancer listens for incoming requests on a specified port.
2. When a request is received, the load balancer selects the server with the least number of active connections and forwards the request to that server.
3. The server processes the request and sends the response back to the load balancer.
4. The load balancer sends the response back to the client.

## Getting Started

1. Clone the repository to your local machine.

2. Install Golang if not already installed.

3. Run the load balancer:
   ```bash
   go build
   ./load-balancer
   ```

# Simple Round-Robin Load Balancer in Go

This is a basic load balancer implemented in Go that distributes incoming requests to a list of servers using the round-robin algorithm.

## How It Works

The load balancer listens for incoming requests and routes them to different backend servers. Here's how it functions:

- Requests are distributed in a round-robin fashion, ensuring each server gets an equal share of the load.

- Health checks can be implemented for each server, but in this code, all servers are considered healthy.

## Getting Started

1. Clone the repository to your local machine.

2. Install Golang if not already installed.

3. Run the load balancer:
   ```bash
   go run main.go
   ```

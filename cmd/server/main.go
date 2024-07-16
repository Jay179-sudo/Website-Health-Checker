package main

import (
	"flag"
	"fmt"
	"jaypd/healthcheck/rpc"
	service "jaypd/healthcheck/url-service"
	"net"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 50050, "The main gRPC process runs on port 500050 by default")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)

	URLService := service.NewURLService()
	server := grpc.NewServer()
	rpc.RegisterURLServiceServer(server, URLService)

	conn, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Could not start TCP Listener. Error: %v", err.Error())
		return
	}
	fmt.Printf("Starting server on port %s\n", address)
	err = server.Serve(conn)
	if err != nil {
		fmt.Printf("Could not start gRPC service. Error: %v", err.Error())
	}

}

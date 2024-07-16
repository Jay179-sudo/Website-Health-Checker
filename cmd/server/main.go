package main

import (
	"flag"
	"fmt"
	"jaypd/healthcheck/rpc"
	service "jaypd/healthcheck/url-service"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
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
		logger.Error("Server", "Error", err.Error())
		return
	}
	logger.Info("Server", "Started at address", address)
	err = server.Serve(conn)
	if err != nil {
		logger.Error("Server", "Error", err.Error())
		return
	}

}

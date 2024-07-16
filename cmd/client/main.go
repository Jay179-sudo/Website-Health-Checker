package main

import (
	"context"
	"flag"
	"fmt"
	"jaypd/healthcheck/rpc"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func main() {
	address := flag.String("address", "0.0.0.0:8080", "default port where the client reaches out to the server")
	url := flag.String("url", "https://google.com", "The URL which the gRPC client sends the request to")
	flag.Parse()

	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Client", "Error", err.Error())
		return
	}
	client := rpc.NewURLServiceClient(conn)

	req := &rpc.URL{
		Url: *url,
	}
	resp, err := client.GetHealthResponse(context.Background(), req)
	if err != nil {
		logger.Error("Client", "Error", err.Error())
		return
	}
	logger.Info("Client", "URL", url, "Response", resp.GetMessage())
	fmt.Printf("%v\n", resp.GetMessage())

}

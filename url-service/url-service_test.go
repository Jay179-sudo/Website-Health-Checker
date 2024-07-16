package service

import (
	"context"
	"jaypd/healthcheck/rpc"
	"net"
	"testing"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
)

func TestURLServiceReturnsFalseForInvalidURLTrueForValid(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	// setup gRPC server
	URLService := NewURLService()
	serverOptions := []grpc.ServerOption{}
	serverOptions = append(serverOptions, grpc.Creds(insecure.NewCredentials()))
	server := grpc.NewServer(serverOptions...)
	rpc.RegisterURLServiceServer(server, URLService)
	go func() {
		if err := server.Serve(lis); err != nil {
			fmt.Printf("Error in the server. Error %v", err.Error())
			return
		}
	}()

	t.Cleanup(func() {
		server.Stop()
	})

	// setup gRPC Client
	resolver.SetDefaultScheme("passthrough")
	conn, err := grpc.NewClient("bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	client := rpc.NewURLServiceClient(conn)

	cases := []struct {
		Name         string
		URL          string
		expectsError bool
	}{
		{
			Name:         "Invalid URL Provided",
			URL:          "wh@t",
			expectsError: true,
		},
		{
			Name:         "Valid URL Provided",
			URL:          "https://google.com",
			expectsError: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			//send request
			req := rpc.URL{
				Url: tt.URL,
			}

			_, err := client.GetHealthResponse(context.Background(), &req)
			if tt.expectsError == false && err != nil {
				t.Errorf("Expected no errors. Received %v", err)
			} else if tt.expectsError == true && err == nil {
				t.Errorf("Expected errors, received %v", err)
			}
		})
	}

}

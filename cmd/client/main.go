package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"jaypd/healthcheck/rpc"
	"log"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load the certificate of the CA who signed the server's certificate
	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}
	// create creds and return
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil

}

func main() {
	address := flag.String("address", "0.0.0.0:8080", "default port where the client reaches out to the server")
	enableTLS := flag.Bool("tls", true, "enable SSL/TLS")
	url := flag.String("url", "https://google.com", "The URL which the gRPC client sends the request to")
	flag.Parse()
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if *enableTLS {
		tlsCreds, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials ", err)
		}
		transportOption = grpc.WithTransportCredentials(tlsCreds)
	}
	conn, err := grpc.NewClient(*address, transportOption)
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

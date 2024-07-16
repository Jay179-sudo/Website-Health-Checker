package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"jaypd/healthcheck/rpc"
	service "jaypd/healthcheck/url-service"
	"log"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/cert/server-cert.pem", "/cert/server-decrypted-key.pem")
	if err != nil {
		return nil, err
	}
	pemClientCA, err := os.ReadFile("/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// create config and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	port := flag.Int("port", 50050, "The main gRPC process runs on port 500050 by default")
	enableTLS := flag.Bool("tls", true, "enable SSL/TLS")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)

	// load TLS configurations
	serverOptions := []grpc.ServerOption{}
	if *enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatalf("Could not load TLS credentials %v", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	URLService := service.NewURLService()
	server := grpc.NewServer(serverOptions...)
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

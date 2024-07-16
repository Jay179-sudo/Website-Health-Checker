package service

import (
	"context"
	"fmt"
	"jaypd/healthcheck/rpc"
)

type URLService struct {
	rpc.UnimplementedURLServiceServer
}

func (u *URLService) GetHealthResponse(ctx context.Context, url *rpc.URL) (*rpc.URLResponse, error) {
	fmt.Print("Received a request")
	return nil, nil
}

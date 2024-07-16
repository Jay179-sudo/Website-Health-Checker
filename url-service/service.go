package service

import (
	"context"
	"errors"
	"fmt"
	"jaypd/healthcheck/rpc"
	"net/http"
	"time"
)

const (
	CONTEXT_TIMEOUT = "the request has timed out"
)

type URLService struct {
	rpc.UnimplementedURLServiceServer
}

func NewURLService() *URLService {
	return &URLService{}
}
func getResponse(ctx context.Context, ch chan string, url string) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		ch <- "Unable to create a HTTP request"
	}

	resp, err := http.DefaultClient.Do(req)
	// TODO perform a more robust check on the status code
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		ch <- fmt.Sprintf("Could not process the URL. The GET request returned %v", resp.StatusCode)
		return
	}

	if err != nil {
		ch <- "There was an error processing the URL"
		return
	}

	ch <- fmt.Sprintf("The %v URL has been processed in %v seconds", url, time.Since(start))

}
func (u *URLService) GetHealthResponse(ctx context.Context, url *rpc.URL) (*rpc.URLResponse, error) {
	// apply a 10 second timeout
	// TODO: add an option to change the number of times a request will be made
	// TODO add structured logging using log/slog
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ch := make(chan string, 1)
	go getResponse(ctx, ch, url.Url)
	var returnMessage string
	resp := &rpc.URLResponse{}

	select {
	case returnMessage = <-ch:
		resp.Message = returnMessage
	case <-ctx.Done():
		return nil, errors.New(CONTEXT_TIMEOUT)
	}
	return resp, nil
}

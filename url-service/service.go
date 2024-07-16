package service

import (
	"context"
	"fmt"
	"jaypd/healthcheck/rpc"
	"net/http"
	"net/url"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	URL_ERROR       = "the server could not process the url"
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
		ch <- URL_ERROR
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- "There was an error processing the URL"
		return
	}
	// TODO perform a more robust check on the status code
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		ch <- URL_ERROR
		return
	}

	ch <- fmt.Sprintf("The %v URL has been processed in %v seconds", url, time.Since(start))

}
func (u *URLService) GetHealthResponse(ctx context.Context, ur *rpc.URL) (*rpc.URLResponse, error) {
	// apply a 10 second timeout
	fmt.Println("Received a URL Request")
	parsedUrl, err := url.ParseRequestURI(ur.Url)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, URL_ERROR)
	}
	// TODO: add an option to change the number of times a request will be made
	// TODO add structured logging using log/slog
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ch := make(chan string, 1)
	go getResponse(ctx, ch, parsedUrl.String())
	var returnMessage string
	resp := &rpc.URLResponse{}

	select {
	case returnMessage = <-ch:
		if returnMessage == URL_ERROR {
			return nil, status.Error(codes.Internal, URL_ERROR)
		}
		resp.Message = returnMessage
		return resp, nil
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, CONTEXT_TIMEOUT)
	}
}

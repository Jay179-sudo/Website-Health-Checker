package service

import (
	"context"
	"fmt"
	"jaypd/healthcheck/rpc"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	URL_ERROR       = "the server could not process the url"
	CONTEXT_TIMEOUT = "the request has timed out"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
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
		ch <- URL_ERROR
		return
	}
	// TODO perform a more robust check on the status code
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		ch <- URL_ERROR
		return
	}

	ch <- fmt.Sprintf("The %v URL has been processed in %v", url, time.Since(start))

}
func (u *URLService) GetHealthResponse(ctx context.Context, ur *rpc.URL) (*rpc.URLResponse, error) {
	// apply a 10 second timeout
	logger.Info("/URLService/GetHealthResponse", "URL:", ur.Url)
	parsedUrl, err := url.ParseRequestURI(ur.Url)
	if err != nil {
		logger.Error("/URLService/GetHealthResponse", "URL", ur.Url, "Code", codes.InvalidArgument, "Error Message", URL_ERROR)
		return nil, status.Error(codes.InvalidArgument, URL_ERROR)
	}
	// TODO: add an option to change the number of times a request will be made
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ch := make(chan string, 1)
	go getResponse(ctx, ch, parsedUrl.String())
	var returnMessage string
	resp := &rpc.URLResponse{}

	select {
	case returnMessage = <-ch:
		if returnMessage == URL_ERROR {
			logger.Error("/URLService/GetHealthResponse", "URL", ur.Url, "Code", codes.Internal, "Error Message", URL_ERROR)
			return nil, status.Error(codes.Internal, URL_ERROR)
		}
		resp.Message = returnMessage
		logger.Info("/URLService/GetHealthResponse", "URL", ur.Url, "Message", returnMessage)
		return resp, nil
	case <-ctx.Done():
		logger.Error("/URLService/GetHealthResponse", "URL", ur.Url, "Code", codes.DeadlineExceeded, "Error Message", CONTEXT_TIMEOUT)
		return nil, status.Error(codes.DeadlineExceeded, CONTEXT_TIMEOUT)
	}
}

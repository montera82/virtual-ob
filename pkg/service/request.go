package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"virtual-orb/pkg/domain"
)

// request represents a struct that holds properties
// required for executing HTTP requests.
type (
	request struct {
		baseURL string
		client  domain.HttpClient
		cb      domain.CircuitBreaker
	}
)

// NewRequestSvc creates a new instance of the request service.
// It requires a base URL, an HTTP client and a circuit breaker.
//
// baseURL: The base URL to which the HTTP requests will be sent.
// client: The HTTP client that will be used to send requests.
// cb: The circuit breaker that will be used to handle failures.
//
// Returns a pointer to a request service instance.
func NewRequestSvc(baseURL string, client domain.HttpClient, cb domain.CircuitBreaker) *request {
	return &request{
		baseURL,
		client,
		cb,
	}
}

// Post sends a POST request to the given route with the provided body.
//
// route: The endpoint route to send the POST request to (relative to baseURL).
// body: The payload/body for the POST request.
//
// Returns the HTTP status code of the response and an error if any occurred during the process.
func (r *request) Post(route string, body any) (httpStatus int, err error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Post: %w", domain.ErrMarshallingPayload)
	}

	var statusCode int

	action := func() (any, error) {
		url := fmt.Sprintf("%s%s", r.baseURL, route)
		req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

		if err != nil {
			return nil, fmt.Errorf("Post: %w", domain.ErrRequestFailed)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Post: %w", domain.ErrRequestFailed)
		}

		defer resp.Body.Close()
		return resp.StatusCode, nil
	}

	result, err := r.cb.Execute(action)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Post: %w", domain.ErrExecutionFailed)
	}

	statusCode, ok := result.(int)
	if !ok {
		return http.StatusInternalServerError, fmt.Errorf("Post: %w", domain.ErrCasting)
	}
	return statusCode, nil
}

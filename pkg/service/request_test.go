package service_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"virtual-orb/mock"
	"virtual-orb/pkg/service"
	testhelper "virtual-orb/test_helper"

	"github.com/sony/gobreaker"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, string, *mock.HttpClient, *mock.CircuitBreaker)
	}{
		{"should successfully post", testSuccessfullyPost},
		{"should handle client error", testHandleClientError},
		{"should handle circuit breaker error", testHandleCBError},
		{"should handle JSON marshal error", testHandleJSONError},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			b := "http://test-base-url/api/v1"
			h := new(mock.HttpClient)
			cb := new(mock.CircuitBreaker)
			test.function(t, b, h, cb)
		})
	}
}

func testSuccessfullyPost(t *testing.T, baseUrl string, h *mock.HttpClient, cb *mock.CircuitBreaker) {
	h.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("OK")),
		}, nil
	}

	cb.ExecuteFunc = func(action any) (any, error) {
		return 200, nil
	}

	r := service.NewRequestSvc(baseUrl, h, cb)
	body := []byte("body")
	_, err := r.Post("/test", body)
	testhelper.Ok(t, err)
	testhelper.Assert(t, cb.ExecuteFuncInvoked == true, "expected httpClient.Do to be invoked")
}

func testHandleClientError(t *testing.T, baseUrl string, h *mock.HttpClient, cb *mock.CircuitBreaker) {
	h.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("client error")
	}
	cbSettings := gobreaker.Settings{}
	cbReal := gobreaker.NewCircuitBreaker(cbSettings)
	r := service.NewRequestSvc(baseUrl, h, cbReal)
	body := []byte("body")
	_, err := r.Post("/test", body)
	testhelper.Assert(t, err != nil, "expected an error from client")
}

func testHandleCBError(t *testing.T, baseUrl string, h *mock.HttpClient, cb *mock.CircuitBreaker) {
	cb.ExecuteFunc = func(action any) (any, error) {
		return nil, errors.New("circuit breaker error")
	}

	r := service.NewRequestSvc(baseUrl, h, cb)
	body := []byte("body")
	_, err := r.Post("/test", body)
	testhelper.Assert(t, err != nil, "expected an error from circuit breaker")
}

func testHandleJSONError(t *testing.T, baseUrl string, h *mock.HttpClient, cb *mock.CircuitBreaker) {
	r := service.NewRequestSvc(baseUrl, h, cb)
	body := make(chan int) // A type that can't be marshaled to JSON
	_, err := r.Post("/test", body)
	testhelper.Assert(t, err != nil, "expected a JSON marshalling error")
}

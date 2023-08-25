package mock

import "net/http"

type (
	HttpClient struct {
		DoFunc        func(req *http.Request) (*http.Response, error)
		DoFuncInvoked bool
	}
)

func (h *HttpClient) Do(req *http.Request) (*http.Response, error) {
	h.DoFuncInvoked = true
	return h.DoFunc(req)
}

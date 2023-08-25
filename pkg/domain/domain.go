package domain

import (
	"net/http"

	"github.com/bwmarrin/snowflake"
)

// Status represents the status of a system, including battery, CPU usage, CPU temperature, and disk space.
type Status struct {
	Battery   float32 `json:"battery"`   // Battery level in percentage.
	CPUUsage  float32 `json:"cpuUsage"`  // CPU usage in percentage.
	CPUTemp   float32 `json:"cpuTemp"`   // CPU temperature in Celsius.
	DiskSpace float32 `json:"diskSpace"` // Available disk space.
}

// Iris represents the iris code and its associated ID.
type Iris struct {
	Id       string `json:"id"`
	IrisCode string `json:"irisCode"`
}

// StatusSvc provides an interface for reporting system status.
type StatusSvc interface {
	// Report takes a status and reports it, returning an error if any.
	Report(s *Status) (err error)
}

// SignUpSvc provides an interface for handling user sign-up based on iris data.
type SignUpSvc interface {
	// SignUp processes the given image and signs it up, returning an error if any.
	SignUp(img []byte) (err error)
}

// RequestSvc provides an interface for making HTTP POST requests.
type RequestSvc interface {
	// Post sends a POST request to the given path with the provided body, returning an HTTP status and an error if any.
	Post(path string, body any) (httpStatus int, err error)
}

// HttpClient is an interface representing the capability to execute HTTP requests.
type HttpClient interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(req *http.Request) (*http.Response, error)
}

// CircuitBreaker is an interface providing the capability to execute functions with circuit breaker logic.
type CircuitBreaker interface {
	// Execute performs the given function with circuit breaker logic and returns the result and an error if any.
	Execute(func() (any, error)) (any, error)
}

// SnowFlakeNode provides an interface for generating unique Snowflake IDs.
type SnowFlakeNode interface {
	// Generate returns a unique snowflake ID.
	Generate() snowflake.ID
}

// SystemInfo is an interface representing the capability to fetch system-related information.
type SystemInfo interface {
	// GetSystemInfo returns a struct containing information about the system's status.
	GetSystemInfo() *Status
}

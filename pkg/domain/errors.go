package domain

import "errors"

var (
	ErrDecode             = errors.New("image decoding failed")
	ErrImageHash          = errors.New("image hashing failed")
	ErrRequestFailed      = errors.New("request failed")
	ErrInvalidImageFormat = errors.New("provided image is not a PNG")
	ErrMarshallingPayload = errors.New("marshalling payload failed")
	ErrCasting            = errors.New("casting failed")
	ErrExecutionFailed    = errors.New("circuit breaker execution failed")
)

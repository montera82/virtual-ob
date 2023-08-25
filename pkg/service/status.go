package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"virtual-orb/pkg/domain"
)

// statusSvc provides services related to reporting the status of the system.
type statusSvc struct {
	requestSvc domain.RequestSvc
	systemInfo domain.SystemInfo
}

// NewStatusSvc initializes a new instance of statusSvc.
//
// requestSvc: Service used to handle HTTP requests.
// systemInfo: Entity responsible for retrieving system-related information.
//
// Returns a pointer to an initialized statusSvc instance.
func NewStatusSvc(requestSvc domain.RequestSvc, systemInfo domain.SystemInfo) *statusSvc {
	return &statusSvc{
		requestSvc,
		systemInfo,
	}
}

// Report gathers system information and reports it.
// The system information is marshaled into a JSON string and then sent as a payload to a "/status" endpoint.
//
// Returns an error if any occurred during the process.
func (ss *statusSvc) Report() error {
	status := ss.systemInfo.GetSystemInfo()

	payload, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("Report: %w", domain.ErrMarshallingPayload)
	}
	payloadJson := string(payload)

	statusCode, err := ss.requestSvc.Post("/status", payloadJson)
	if err != nil || statusCode != http.StatusOK {
		return fmt.Errorf("Report: %w", domain.ErrRequestFailed)
	}

	return nil
}

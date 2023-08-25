package service_test

import (
	"errors"
	"testing"
	"virtual-orb/mock"
	"virtual-orb/pkg/domain"
	"virtual-orb/pkg/service"
	testhelper "virtual-orb/test_helper"
)

func TestStatusService(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *mock.RequestSvc, *mock.SystemInfo)
	}{
		{"should report status successfully", testSuccessfulStatusReport},
		{"should handle post request error", testStatusPostRequestError},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			reqSvc := new(mock.RequestSvc)
			sysInfo := new(mock.SystemInfo)
			test.function(t, reqSvc, sysInfo)
		})
	}
}

func testSuccessfulStatusReport(t *testing.T, reqSvc *mock.RequestSvc, sysInfo *mock.SystemInfo) {
	sysInfo.GetSystemInfoFunc = func() *domain.Status {
		return &domain.Status{}
	}
	reqSvc.PostFunc = func(path string, body any) (httpStatus int, err error) {
		return 200, nil
	}
	statusService := service.NewStatusSvc(reqSvc, sysInfo)
	err := statusService.Report()
	testhelper.Ok(t, err)
}

func testStatusPostRequestError(t *testing.T, reqSvc *mock.RequestSvc, sysInfo *mock.SystemInfo) {
	sysInfo.GetSystemInfoFunc = func() *domain.Status {
		return &domain.Status{}
	}
	reqSvc.PostFunc = func(path string, body any) (httpStatus int, err error) {
		return 500, errors.New("mocked post error")
	}
	statusService := service.NewStatusSvc(reqSvc, sysInfo)
	err := statusService.Report()
	testhelper.Assert(t, err != nil, "expected a post request error")
}

package service_test

import (
	"errors"
	"fmt"
	"testing"

	"virtual-orb/mock"
	"virtual-orb/pkg/domain"
	"virtual-orb/pkg/platform"
	"virtual-orb/pkg/service"
	testhelper "virtual-orb/test_helper"

	"github.com/bwmarrin/snowflake"
)

func TestSignUpService(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *mock.RequestSvc, *mock.SnowFlakeNode)
	}{
		{"should handle and sign image properly", testHandleAndSignImage},
		{"should handle image decoding error", testImageDecodingError},
		{"should reject non-PNG image format", testNonPNGImage},
		{"should handle post request error", testPostRequestError},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			reqSvc := new(mock.RequestSvc)
			sfNode := new(mock.SnowFlakeNode)
			test.function(t, reqSvc, sfNode)
		})
	}
}

func testHandleAndSignImage(t *testing.T, reqSvc *mock.RequestSvc, sfNode *mock.SnowFlakeNode) {
	img, _ := platform.GenerateRandomImageData()
	reqSvc.PostFunc = func(path string, body any) (httpStatus int, err error) {
		return 201, nil
	}
	sfNode.GenerateFunc = func() snowflake.ID {
		return snowflake.ID(123456789)
	}
	signKey := "test-key"
	signUpService := service.NewSignUpSvc(signKey, sfNode, reqSvc)
	err := signUpService.SignUp(img)
	testhelper.Ok(t, err)
}

func testImageDecodingError(t *testing.T, reqSvc *mock.RequestSvc, sfNode *mock.SnowFlakeNode) {
	badImg := []byte("not a valid image")
	signUpService := service.NewSignUpSvc("test-key", sfNode, reqSvc)
	err := signUpService.SignUp(badImg)
	testhelper.Assert(t, err != nil, "expected an image decoding error")
}

func testNonPNGImage(t *testing.T, reqSvc *mock.RequestSvc, sfNode *mock.SnowFlakeNode) {
	nonPngImg := []byte("GIF89a...")
	signUpService := service.NewSignUpSvc("test-key", sfNode, reqSvc)
	err := signUpService.SignUp(nonPngImg)
	testhelper.Assert(t, err != nil && errors.Is(err, domain.ErrInvalidImageFormat), "expected an error for non-PNG image format")
}

func testPostRequestError(t *testing.T, reqSvc *mock.RequestSvc, sfNode *mock.SnowFlakeNode) {
	img, _ := platform.GenerateRandomImageData()
	reqSvc.PostFunc = func(path string, body any) (httpStatus int, err error) {
		return 500, fmt.Errorf("Post request failed")
	}
	sfNode.GenerateFunc = func() snowflake.ID {
		return snowflake.ID(123456789)
	}
	signUpService := service.NewSignUpSvc("test-key", sfNode, reqSvc)
	err := signUpService.SignUp(img)
	testhelper.Assert(t, err != nil, "expected a post request error")
}

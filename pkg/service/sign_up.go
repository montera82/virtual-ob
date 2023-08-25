package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"image/png"
	"net/http"
	"virtual-orb/pkg/domain"

	"github.com/corona10/goimagehash"
)

// signUpSvc encapsulates services required for user sign-up,
// especially with an emphasis on image processing and cryptographic security.
type (
	signUpSvc struct {
		signKey       string
		snowflakeNode domain.SnowFlakeNode
		requestSvc    domain.RequestSvc
	}
)

// NewSignUpSvc initializes a new signUpSvc instance.
//
// signKey: Secret key used for signing operations.
// snowflakeNode: Entity responsible for generating unique IDs.
// requestSvc: Service to handle HTTP requests.
//
// Returns a pointer to an initialized signUpSvc instance.
func NewSignUpSvc(signKey string, snowflakeNode domain.SnowFlakeNode, requestSvc domain.RequestSvc) *signUpSvc {
	return &signUpSvc{
		signKey,
		snowflakeNode,
		requestSvc,
	}
}

// SignUp processes a user sign-up request using an image (iris scan).
// The image is perceptually hashed, signed, and sent for further processing.
//
// img: The image data in bytes.
//
// Returns an error if any occurred during the process.
func (s *signUpSvc) SignUp(img []byte) error {

	imgType := http.DetectContentType(img)
	if imgType != "image/png" {
		return fmt.Errorf("SignUp: %w", domain.ErrInvalidImageFormat)
	}

	i, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return fmt.Errorf("SignUp: %w", domain.ErrDecode)
	}

	// Hash the image to generate iris code.
	// Calculate the iris code locally to maximize privacy.
	// References:
	// https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
	// Todo read https://tech.okcupid.com/evaluating-perceptual-image-hashes-at-okcupid-e98a3e74aa3a
	irisCode, err := goimagehash.AverageHash(i)
	if err != nil {
		return fmt.Errorf("SignUp: %w", domain.ErrImageHash)
	}

	// Sign the iris code for security verification.
	signedIrisCode := s.signIrisCode(irisCode.ToString())

	id := s.snowflakeNode.Generate().String()
	request := domain.Iris{
		Id:       id,
		IrisCode: signedIrisCode,
	}

	statusCode, err := s.requestSvc.Post("/sign-up", request)
	if err != nil || statusCode != http.StatusCreated {
		return fmt.Errorf("SignUp: %w", domain.ErrRequestFailed)
	}
	return nil
}

// signIrisCode signs the provided iris code using HMAC-SHA256 and the service's signKey.
//
// irisCode: The perceptual hash of the iris scan.
//
// Returns the signed iris code in hex string format.
func (s *signUpSvc) signIrisCode(irisCode string) string {
	mac := hmac.New(sha256.New, []byte(s.signKey))
	mac.Write([]byte(irisCode))
	return hex.EncodeToString(mac.Sum(nil))
}

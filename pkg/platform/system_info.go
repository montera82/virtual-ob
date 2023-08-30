// Package platform provides utilities related to system information.
package platform

import (
	"math/rand"
	"time"
	"virtual-orb/pkg/domain"
)

type (
	// systemInfo represents an implementation of the SystemInfo interface
	// from the domain package. It provides mock system status information.
	systemInfo struct {
	}
)

// init initializes the package by seeding the random number generator
// using the current time.
func init() {
	rand.NewSource(time.Now().UnixNano())
}

// NewSystemInfo creates a new instance of systemInfo which implements
// the domain.SystemInfo interface. It provides mock system status details
func NewSystemInfo() domain.SystemInfo {
	return &systemInfo{}
}

// GetSystemInfo provides mock system status details, generating random values
// for battery percentage, CPU usage, CPU temperature, and available disk space.
// This method satisfies the SystemInfo interface of the domain package.
func (s *systemInfo) GetSystemInfo() *domain.Status {

	status := &domain.Status{
		Battery:  float32(rand.Intn(101)),
		CPUUsage: float32(rand.Intn(101)),
		// Assume realistic temperatures between -30 and 60
		CPUTemp:   -30 + 90*rand.Float32(),
		DiskSpace: float32(rand.Intn(501)),
	}

	return status
}

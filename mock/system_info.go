package mock

import "virtual-orb/pkg/domain"

type (
	SystemInfo struct {
		GetSystemInfoFunc func() *domain.Status
	}
)

func (m *SystemInfo) GetSystemInfo() *domain.Status {
	return m.GetSystemInfoFunc()
}

package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, fmt.Errorf("mock error: %w", args.Error(1))
	}

	if ifaces, ok := args.Get(0).([]*wifi.Interface); ok {
		return ifaces, fmt.Errorf("mock error: %w", args.Error(1))
	}

	return nil, fmt.Errorf("mock error: %w", args.Error(1))
}

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

	var result []*wifi.Interface
	if args.Get(0) != nil {
		result = args.Get(0).([]*wifi.Interface)
	}

	return result, args.Error(1)
}

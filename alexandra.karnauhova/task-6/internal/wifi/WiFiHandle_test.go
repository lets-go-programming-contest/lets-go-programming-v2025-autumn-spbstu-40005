package wifi_test

import (
	"fmt"

	wifi "github.com/mdlayher/wifi"
	mock "github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Interfaces")
	}

	if rf, ok := ret.Get(0).(func() ([]*wifi.Interface, error)); ok {
		return rf()
	}

	var r0 []*wifi.Interface
	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*wifi.Interface)
	}

	r1 := ret.Error(1)
	if r1 != nil {
		r1 = fmt.Errorf("mock Interfaces error: %w", r1)
	}

	return r0, r1
}

func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(fn func())
}) *WiFiHandle {
	mock := &WiFiHandle{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

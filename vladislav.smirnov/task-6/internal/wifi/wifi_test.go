package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	myWifi "github.com/smirnov-vladislav/task-6/internal/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --name=WiFiHandle --dir=../wifi --output=. --outpkg=wifi_test --testonly --quiet

var errTest = errors.New("test error")

func makeMAC(addr string) net.HardwareAddr {
	mac, err := net.ParseMAC(addr)
	if err != nil {
		panic("invalid MAC address: " + addr)
	}
	return mac
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	assert.NotNil(t, service)
	assert.Equal(t, mockHandler, service.WiFi)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: makeMAC("00:11:22:33:44:55")},
		{Name: "wlan1", HardwareAddr: makeMAC("AA:BB:CC:DD:EE:FF")},
	}
	mockHandler.On("Interfaces").Return(interfaces, nil).Once()

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, names)
	mockHandler.AssertExpectations(t)
}

func TestGetNames_EmptyList(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Empty(t, names)
	mockHandler.AssertExpectations(t)
}

func TestGetNames_Failure(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return(nil, errTest).Once()

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.ErrorContains(t, err, "getting interfaces")
	mockHandler.AssertExpectations(t)
}

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	interfaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: makeMAC("11:22:33:44:55:66")},
		{Name: "wifi1", HardwareAddr: makeMAC("AA:BB:CC:DD:EE:FF")},
	}
	mockHandler.On("Interfaces").Return(interfaces, nil).Once()

	macs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		makeMAC("11:22:33:44:55:66"),
		makeMAC("AA:BB:CC:DD:EE:FF"),
	}, macs)
	mockHandler.AssertExpectations(t)
}

func TestGetAddresses_WithNilAddress(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	interfaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: nil},
		{Name: "wifi1", HardwareAddr: makeMAC("11:22:33:44:55:66")},
	}
	mockHandler.On("Interfaces").Return(interfaces, nil).Once()

	macs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, macs, 2)
	require.Nil(t, macs[0])
	require.Equal(t, makeMAC("11:22:33:44:55:66"), macs[1])
	mockHandler.AssertExpectations(t)
}

func TestGetAddresses_EmptyResult(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

	macs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Empty(t, macs)
	mockHandler.AssertExpectations(t)
}

func TestGetAddresses_Failure(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return(nil, errTest).Once()

	macs, err := service.GetAddresses()

	require.Error(t, err)
	require.Nil(t, macs)
	require.ErrorContains(t, err, "getting interfaces")
	mockHandler.AssertExpectations(t)
}

func TestMultipleCalls(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: makeMAC("00:11:22:33:44:55")},
	}

	mockHandler.On("Interfaces").Return(interfaces, nil).Twice()

	names1, err1 := service.GetNames()
	require.NoError(t, err1)
	require.Equal(t, []string{"wlan0"}, names1)

	names2, err2 := service.GetNames()
	require.NoError(t, err2)
	require.Equal(t, []string{"wlan0"}, names2)

	mockHandler.AssertExpectations(t)
}

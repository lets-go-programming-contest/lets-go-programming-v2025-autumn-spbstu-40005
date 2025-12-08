package wifi

import (
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"alina.duhanina/task-6/internal/wifi"
)

func TestWiFiService_GetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr1, _ := net.ParseMAC("00:11:22:33:44:55")
	hwAddr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: hwAddr1},
		{Name: "wlan1", HardwareAddr: hwAddr2},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Len(t, addrs, 2)
	assert.Equal(t, hwAddr1, addrs[0])
	assert.Equal(t, hwAddr2, addrs[1])

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	interfaces := []*wifi.Interface{}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Empty(t, addrs)

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, fmt.Errorf("mock error"))

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	assert.Error(t, err)
	assert.Nil(t, addrs)
	assert.Contains(t, err.Error(), "getting interfaces")

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr, _ := net.ParseMAC("00:11:22:33:44:55")
	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: hwAddr},
		{Name: "wlan1"},
		{Name: "eth0"},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	interfaces := []*wifi.Interface{}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Empty(t, names)

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, fmt.Errorf("mock error"))

	service := New(mockWiFi)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "getting interfaces")

	mockWiFi.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	service := New(mockWiFi)
	assert.NotNil(t, service)
	assert.Equal(t, mockWiFi, service.WiFi)
}

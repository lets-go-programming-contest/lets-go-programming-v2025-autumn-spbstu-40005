package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "kirill.kharlamov/task-6/internal/wifi"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var ErrTest = errors.New("test error")

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	ifaces := []*wifi.Interface{
		{Name: "eth0", HardwareAddr: parseMAC(t, "11:22:33:44:55:66")},
		{Name: "wlan0", HardwareAddr: parseMAC(t, "AA:BB:CC:DD:EE:FF")},
	}
	mockWiFi.On("Interfaces").Return(ifaces, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		parseMAC(t, "11:22:33:44:55:66"),
		parseMAC(t, "AA:BB:CC:DD:EE:FF"),
	}, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return(nil, ErrTest)

	addrs, err := service.GetAddresses()
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, addrs)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	ifaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "ethernet1", HardwareAddr: parseMAC(t, "66:77:88:99:AA:BB")},
	}
	mockWiFi.On("Interfaces").Return(ifaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wifi0", "ethernet1"}, names)
}

func TestGetNames_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return(nil, ErrTest)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, names)
}

func parseMAC(t *testing.T, s string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(s)
	require.NoError(t, err)

	return hwAddr
}

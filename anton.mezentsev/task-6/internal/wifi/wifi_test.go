package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifi "anton.mezentsev/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var ErrMock = errors.New("mock error")

func TestGetMACs_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := wifi.NewService(mockWiFi)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: parseMAC(t, "08:00:27:ab:cd:ef")},
		{Name: "wlan1", HardwareAddr: parseMAC(t, "08:00:27:12:34:56")},
	}
	mockWiFi.On("GetInterfaces").Return(ifaces, nil)

	addrs, err := service.GetMACs()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		parseMAC(t, "08:00:27:ab:cd:ef"),
		parseMAC(t, "08:00:27:12:34:56"),
	}, addrs)
}

func TestGetMACs_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := wifi.NewService(mockWiFi)

	mockWiFi.On("GetInterfaces").Return(nil, ErrMock)

	addrs, err := service.GetMACs()
	require.ErrorContains(t, err, "failed to get interfaces")
	require.Nil(t, addrs)
}

func TestGetInterfaceNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := wifi.NewService(mockWiFi)

	ifaces := []*wifi.Interface{
		{Name: "wlp2s0", HardwareAddr: parseMAC(t, "aa:bb:cc:dd:ee:ff")},
		{Name: "eth0", HardwareAddr: parseMAC(t, "11:22:33:44:55:66")},
	}
	mockWiFi.On("GetInterfaces").Return(ifaces, nil)

	names, err := service.GetInterfaceNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlp2s0", "eth0"}, names)
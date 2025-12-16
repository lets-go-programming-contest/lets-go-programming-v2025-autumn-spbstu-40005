package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "polina.gavrilova/task-6/internal/wifi"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrDriverNotLoaded  = errors.New("driver not loaded")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &WiFiHandle{}
		service := myWifi.New(mockWiFi)

		hw1, _ := net.ParseMAC("01:23:45:67:89:00")
		hw2, _ := net.ParseMAC("ab:cd:ef:hg:ik:lm")

		mockWiFi.On("Interfaces").Return([]*wifi.Interface{
			{Name: "wlan0", HardwareAddr: hw1},
			{Name: "eth0", HardwareAddr: hw2},
		}, nil)

		addrs, err := service.GetAddresses()
		require.NoError(t, err)
		require.Equal(t, []net.HardwareAddr{hw1, hw2}, addrs)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &WiFiHandle{}
		service := myWifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, ErrPermissionDenied)

		addrs, err := service.GetAddresses()
		require.Error(t, err)
		require.Nil(t, addrs)
		require.Contains(t, err.Error(), "getting interfaces:")
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &WiFiHandle{}
		service := myWifi.New(mockWiFi)

		hw, _ := net.ParseMAC("ab:cd:ef:hg:ik:lm")
		mockWiFi.On("Interfaces").Return([]*wifi.Interface{
			{Name: "wlan0", HardwareAddr: hw},
			{Name: "eth1", HardwareAddr: hw},
		}, nil)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"wlan0", "eth1"}, names)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &WiFiHandle{}
		service := myWifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, ErrDriverNotLoaded)

		names, err := service.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "getting interfaces:")
	})
}

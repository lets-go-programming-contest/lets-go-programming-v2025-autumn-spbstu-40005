package wifi_test

import (
	"errors"
	"net"
	"testing"

	myWifi "github.com/TvoyBatyA1234/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var TestErr = errors.New("test error")

func createMAC(s string) net.HardwareAddr {
	addr, _ := net.ParseMAC(s)

	return addr
}

func TestRetrieveMACs_Successful(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	networkIfaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: createMAC("11:22:33:44:55:66")},
		{Name: "wifi1", HardwareAddr: createMAC("aa:bb:cc:dd:ee:ff")},
	}
	mockHandler.On("Interfaces").Return(networkIfaces, nil)

	macs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		createMAC("11:22:33:44:55:66"),
		createMAC("aa:bb:cc:dd:ee:ff"),
	}, macs)
}

func TestRetrieveMACs_Failed(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return(nil, TestErr)

	macs, err := service.GetAddresses()
	require.Error(t, err)
	require.Contains(t, err.Error(), "getting interfaces")
	require.ErrorIs(t, err, TestErr)
	require.Nil(t, macs)
}

func TestRetrieveInterfaceNames_Successful(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	networkIfaces := []*wifi.Interface{
		{Name: "wireless0", HardwareAddr: createMAC("11:22:33:44:55:66")},
		{Name: "ethernet1", HardwareAddr: createMAC("aa:bb:cc:dd:ee:ff")},
	}
	mockHandler.On("Interfaces").Return(networkIfaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wireless0", "ethernet1"}, names)
}

func TestRetrieveInterfaceNames_Failed(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return(nil, TestErr)

	names, err := service.GetNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "getting interfaces")
	require.ErrorIs(t, err, TestErr)
	require.Nil(t, names)
}

func TestRetrieveMACs_EmptyResult(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	service := myWifi.New(mockHandler)

	mockHandler.On("Interfaces").Return([]*wifi.Interface{}, nil)

	macs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Empty(t, macs)
}

package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "nikita.brevnov/task-6/internal/wifi"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var ErrTest = errors.New("test error")

func TestFetchAddresses_Success(t *testing.T) {
	t.Parallel()

	mockClient := NewWiFiHandle(t)
	service := myWifi.New(mockClient)

	ifaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: parseMAC("11:22:33:44:55:66")},
		{Name: "wifi1", HardwareAddr: parseMAC("aa:bb:cc:dd:ee:ff")},
	}
	mockClient.On("Interfaces").Return(ifaces, nil)

	result, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		parseMAC("11:22:33:44:55:66"),
		parseMAC("aa:bb:cc:dd:ee:ff"),
	}, result)
}

func TestFetchAddresses_Failure(t *testing.T) {
	t.Parallel()

	mockClient := NewWiFiHandle(t)
	service := myWifi.New(mockClient)

	mockClient.On("Interfaces").Return(nil, ErrTest)

	addrs, err := service.GetAddresses()
	require.Error(t, err)
	require.Contains(t, err.Error(), "interfaces")
	require.ErrorIs(t, err, ErrTest)
	require.Nil(t, addrs)
}

func TestFetchInterfaceNames_Success(t *testing.T) {
	t.Parallel()

	mockClient := NewWiFiHandle(t)
	service := myWifi.New(mockClient)

	ifaces := []*wifi.Interface{
		{Name: "wifi0", HardwareAddr: parseMAC("11:22:33:44:55:66")},
		{Name: "eth0", HardwareAddr: parseMAC("22:33:44:55:66:77")},
	}
	mockClient.On("Interfaces").Return(ifaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wifi0", "eth0"}, names)
}

func TestFetchInterfaceNames_Failure(t *testing.T) {
	t.Parallel()

	mockClient := NewWiFiHandle(t)
	service := myWifi.New(mockClient)

	mockClient.On("Interfaces").Return(nil, ErrTest)

	result, err := service.GetNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "interfaces")
	require.ErrorIs(t, err, ErrTest)
	require.Nil(t, result)
}

func parseMAC(s string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(s)
	if err != nil {
		return nil
	}
	return hwAddr
}

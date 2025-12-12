package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "gordey.shapkov/task-6/internal/wifi"
)

var testTable = []struct { //nolint:gochecknoglobals
	addrs       []string
	names       []string
	errWrap     string
	errExpected error
}{
	{
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		names: []string{"eth1", "eth2"},
	},
	{
		addrs:       nil,
		errWrap:     "getting interfaces: ",
		errExpected: ErrExpected,
	},
	{
		addrs: []string{},
		names: []string{},
	},
}

var ErrExpected = errors.New("expected error")

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{WiFi: mockWifi}

	for _, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(helperMockIfaces(t, row.addrs), row.errExpected)

		actualAddrs, err := wifiService.GetAddresses()

		if row.errExpected != nil {
			require.ErrorContains(t, err, row.errWrap)

			continue
		}

		require.NoError(t, err)
		require.Equal(t, helperParseMACs(t, row.addrs), actualAddrs)
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{WiFi: mockWifi}

	for _, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(helperMockIfaces(t, row.addrs), row.errExpected)

		actualNames, err := wifiService.GetNames()

		if row.errExpected != nil {
			require.ErrorContains(t, err, row.errWrap)

			continue
		}

		require.NoError(t, err)
		require.Equal(t, row.names, actualNames)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)
	require.Equal(t, mockWifi, wifiService.WiFi)
}

func helperMockIfaces(t *testing.T, addrs []string) []*wifi.Interface {
	t.Helper()

	interfaces := make([]*wifi.Interface, 0, len(addrs))

	for i, addrStr := range addrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("eth%d", i+1),
			HardwareAddr: hwAddr,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}

		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func helperParseMACs(t *testing.T, macStr []string) []net.HardwareAddr {
	t.Helper()

	addrs := make([]net.HardwareAddr, 0, len(macStr))

	for _, addr := range macStr {
		addrs = append(addrs, parseMAC(addr))
	}

	return addrs
}

func parseMAC(macStr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return nil
	}

	return hwAddr
}

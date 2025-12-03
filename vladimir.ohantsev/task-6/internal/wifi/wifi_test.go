package wifi_test

import (
	"errors"
	"net"
	"testing"

	mywifi "github.com/P3rCh1/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var ErrSome = errors.New("some error")

type testcase struct {
	testName      string
	names         []string
	addrs         []string
	expectedError error
}

func getTests() []testcase {
	return []testcase{
		{
			testName: "success",
			addrs:    []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			names:    []string{"eth1", "eth2"},
		},
		{
			testName:      "error",
			expectedError: ErrSome,
		},
	}
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := mywifi.WiFiService{WiFi: mockWifi}

	for _, test := range getTests() {
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()

			mockWifi.On("Interfaces").Unset()
			mockWifi.On("Interfaces").Return(mockIfaces(&test), test.expectedError)

			actualAddrs, err := wifiService.GetAddresses()

			if test.expectedError != nil {
				require.ErrorIs(
					t, err, test.expectedError,
					"expected: %s, actual: %s", test.expectedError, err,
				)

				return
			}

			require.NoError(
				t, err,
				"error must be nil, %s", err,
			)

			require.Equal(t, parseMACs(test.addrs), actualAddrs,
				"expected: %s, actual: %s", parseMACs(test.addrs), actualAddrs,
			)
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := mywifi.WiFiService{WiFi: mockWifi}

	for _, test := range getTests() {
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()

			mockWifi.On("Interfaces").Unset()
			mockWifi.On("Interfaces").Return(mockIfaces(&test), test.expectedError)

			actualNames, err := wifiService.GetNames()

			if test.expectedError != nil {
				require.ErrorIs(
					t, err, test.expectedError,
					"expected: %s, actual: %s", test.expectedError, err,
				)

				return
			}

			require.NoError(
				t, err,
				"error must be nil, %s", err,
			)

			require.Equal(t, test.names, actualNames,
				"expected: %s, actual: %s", test.names, actualNames,
			)
		})

	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := mywifi.New(mockWifi)
	require.Equal(
		t, mockWifi, wifiService.WiFi,
		"expected: %s, actual: %s", mockWifi, wifiService.WiFi,
	)
}

func mockIfaces(test *testcase) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(test.addrs))

	for i, addr := range test.addrs {
		hwAddr := parseMAC(addr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         test.names[i],
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

func parseMACs(macStr []string) []net.HardwareAddr {
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

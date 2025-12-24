package wifi_test

import (
	"errors"
	"net"
	"testing"

	myWifi "github.com/arseniy.shchadilov/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var errExpected = errors.New("expected error")

type wifiTestCase struct {
	name           string
	mockInterfaces []*wifi.Interface
	mockErr        error
	expectedAddrs  []net.HardwareAddr
	expectedNames  []string
	errContains    string
}

func parseMAC(t *testing.T, s string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(s)
	if err != nil {
		return nil
	}

	return hwAddr
}

//nolint:unparam
func createTestInterface(t *testing.T, name string, mac string) *wifi.Interface {
	t.Helper()

	return &wifi.Interface{
		Name:         name,
		HardwareAddr: parseMAC(t, mac),
		Index:        1,
		PHY:          1,
		Device:       1,
		Type:         wifi.InterfaceTypeAPVLAN,
		Frequency:    0,
	}
}

//nolint:gochecknoglobals
var testCasesGetAddresses = []wifiTestCase{
	{
		name: "success with multiple interfaces",
		mockInterfaces: []*wifi.Interface{
			createTestInterface(nil, "wlan0", "00:11:22:33:44:55"),
			createTestInterface(nil, "eth1", "aa:bb:cc:dd:ee:ff"),
		},
		mockErr: nil,
		expectedAddrs: []net.HardwareAddr{
			parseMAC(nil, "00:11:22:33:44:55"),
			parseMAC(nil, "aa:bb:cc:dd:ee:ff"),
		},
	},
	{
		name:           "success with no interfaces",
		mockInterfaces: []*wifi.Interface{},
		mockErr:        nil,
		expectedAddrs:  []net.HardwareAddr{},
	},
	{
		name:           "interface error",
		mockInterfaces: nil,
		mockErr:        errExpected,
		expectedAddrs:  nil,
		errContains:    "getting interfaces",
	},
	{
		name: "interface with nil MAC address",
		mockInterfaces: []*wifi.Interface{
			{
				Name:         "wlan0",
				HardwareAddr: nil,
				Index:        1,
			},
		},
		mockErr:       nil,
		expectedAddrs: []net.HardwareAddr{nil},
	},
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	for _, tc := range testCasesGetAddresses {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewWiFiHandle(t)
			service := myWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(tc.mockInterfaces, tc.mockErr)

			addrs, err := service.GetAddresses()

			if tc.errContains != "" {
				require.ErrorContains(t, err, tc.errContains)
				require.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedAddrs, addrs)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

//nolint:gochecknoglobals
var testCasesGetNames = []wifiTestCase{
	{
		name: "success with multiple interfaces",
		mockInterfaces: []*wifi.Interface{
			createTestInterface(nil, "wlan0", "00:11:22:33:44:55"),
			createTestInterface(nil, "eth1", "aa:bb:cc:dd:ee:ff"),
			createTestInterface(nil, "lo", "00:00:00:00:00:00"),
		},
		mockErr:       nil,
		expectedNames: []string{"wlan0", "eth1", "lo"},
	},
	{
		name:           "success with no interfaces",
		mockInterfaces: []*wifi.Interface{},
		mockErr:        nil,
		expectedNames:  []string{},
	},
	{
		name:           "interface error",
		mockInterfaces: nil,
		mockErr:        errExpected,
		expectedNames:  nil,
		errContains:    "getting interfaces",
	},
	{
		name: "interface with empty name",
		mockInterfaces: []*wifi.Interface{
			{
				Name:         "",
				HardwareAddr: parseMAC(nil, "00:11:22:33:44:55"),
				Index:        1,
			},
		},
		mockErr:       nil,
		expectedNames: []string{""},
	},
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	for _, tc := range testCasesGetNames {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewWiFiHandle(t)
			service := myWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(tc.mockInterfaces, tc.mockErr)

			names, err := service.GetNames()

			if tc.errContains != "" {
				require.ErrorContains(t, err, tc.errContains)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedNames, names)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	require.NotNil(t, service)
	require.Equal(t, mockWiFi, service.WiFi)
}

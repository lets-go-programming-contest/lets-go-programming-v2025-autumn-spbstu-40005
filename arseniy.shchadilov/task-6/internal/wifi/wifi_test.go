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

func parseMAC(s string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(s)
	if err != nil {
		return nil
	}

	return hwAddr
}

func createTestInterface(name string, mac string) *wifi.Interface {
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: parseMAC(mac),
		Index:        1,
		PHY:          1,
		Device:       1,
		Type:         wifi.InterfaceTypeAPVLAN,
		Frequency:    0,
	}
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	testCases := []wifiTestCase{
		{
			name: "success with multiple interfaces",
			mockInterfaces: []*wifi.Interface{
				createTestInterface("wlan0", "00:11:22:33:44:55"),
				createTestInterface("eth1", "aa:bb:cc:dd:ee:ff"),
			},
			mockErr: nil,
			expectedAddrs: []net.HardwareAddr{
				parseMAC("00:11:22:33:44:55"),
				parseMAC("aa:bb:cc:dd:ee:ff"),
			},
			errContains: "",
		},
		{
			name:           "success with no interfaces",
			mockInterfaces: []*wifi.Interface{},
			mockErr:        nil,
			expectedAddrs:  []net.HardwareAddr{},
			errContains:    "",
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
			errContains:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewWiFiHandle(t)
			service := myWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(tc.mockInterfaces, tc.mockErr)

			addrs, err := service.GetAddresses()

			if tc.mockErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.mockErr)
				require.Contains(t, err.Error(), tc.errContains)
				require.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedAddrs, addrs)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testCases := []wifiTestCase{
		{
			name: "success with multiple interfaces",
			mockInterfaces: []*wifi.Interface{
				createTestInterface("wlan0", "00:11:22:33:44:55"),
				createTestInterface("eth1", "aa:bb:cc:dd:ee:ff"),
				createTestInterface("lo", "00:00:00:00:00:00"),
			},
			mockErr:       nil,
			expectedNames: []string{"wlan0", "eth1", "lo"},
			errContains:   "",
		},
		{
			name:           "success with no interfaces",
			mockInterfaces: []*wifi.Interface{},
			mockErr:        nil,
			expectedNames:  []string{},
			errContains:    "",
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
					HardwareAddr: parseMAC("00:11:22:33:44:55"),
					Index:        1,
				},
			},
			mockErr:       nil,
			expectedNames: []string{""},
			errContains:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewWiFiHandle(t)
			service := myWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(tc.mockInterfaces, tc.mockErr)

			names, err := service.GetNames()

			if tc.mockErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.mockErr)
				require.Contains(t, err.Error(), tc.errContains)
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

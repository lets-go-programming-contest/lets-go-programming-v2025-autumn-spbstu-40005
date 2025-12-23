package wifi_test

import (
	"errors"
	"net"
	"testing"

	myWiFi "alexandra.karnauhova/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var ErrExpected = errors.New("expected error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWiFi.New(mockWiFi)

	require.Equal(t, mockWiFi, service.WiFi)
}

func parseMAC(s string) net.HardwareAddr {
	hwAddr, _ := net.ParseMAC(s)

	return hwAddr
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		mockReturnInterfaces []*wifi.Interface
		mockReturnError      error
		expectedAddrs        []net.HardwareAddr
		expectError          bool
	}{
		{
			name: "success",
			mockReturnInterfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: parseMAC("00:11:22:33:44:55")},
				{Name: "eth0", HardwareAddr: parseMAC("aa:bb:cc:dd:ee:ff")},
			},
			mockReturnError: nil,
			expectedAddrs: []net.HardwareAddr{
				parseMAC("00:11:22:33:44:55"),
				parseMAC("aa:bb:cc:dd:ee:ff"),
			},
			expectError: false,
		},
		{
			name:                 "error from interfaces",
			mockReturnInterfaces: nil,
			mockReturnError:      ErrExpected,
			expectedAddrs:        nil,
			expectError:          true,
		},
		{
			name:                 "empty interfaces",
			mockReturnInterfaces: []*wifi.Interface{},
			mockReturnError:      nil,
			expectedAddrs:        []net.HardwareAddr{},
			expectError:          false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewWiFiHandle(t)
			service := myWiFi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(testCase.mockReturnInterfaces, testCase.mockReturnError)

			addrs, err := service.GetAddresses()

			if testCase.expectError {
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.mockReturnError)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedAddrs, addrs)
			}
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		mockReturnInterfaces []*wifi.Interface
		mockReturnError      error
		expectedNames        []string
		expectError          bool
	}{
		{
			name: "success",
			mockReturnInterfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: parseMAC("00:11:22:33:44:55")},
				{Name: "eth0", HardwareAddr: parseMAC("aa:bb:cc:dd:ee:ff")},
			},
			mockReturnError: nil,
			expectedNames:   []string{"wlan0", "eth0"},
			expectError:     false,
		},
		{
			name:                 "error from interfaces",
			mockReturnInterfaces: nil,
			mockReturnError:      ErrExpected,
			expectedNames:        nil,
			expectError:          true,
		},
		{
			name:                 "empty interfaces",
			mockReturnInterfaces: []*wifi.Interface{},
			mockReturnError:      nil,
			expectedNames:        []string{},
			expectError:          false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewWiFiHandle(t)
			service := myWiFi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(testCase.mockReturnInterfaces, testCase.mockReturnError)

			names, err := service.GetNames()

			if testCase.expectError {
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.mockReturnError)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedNames, names)
			}
		})
	}
}

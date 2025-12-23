package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wifiService "oleg.zholobov/task-6/internal/wifi"
)

var (
	errAccessDenied    = errors.New("access denied")
	errHardwareMissing = errors.New("hardware missing")
)

type wifiTestCase struct {
	testName       string
	mockError      error
	mockInterfaces []*wifi.Interface
	expectError    bool
	expectedResult interface{}
}

func TestWiFiService_GetHardwareAddresses(t *testing.T) {
	t.Parallel()

	testScenarios := []wifiTestCase{
		{
			testName: "successful address retrieval",
			mockInterfaces: []*wifi.Interface{
				createInterface(t, "wlp2s0", "11:22:33:44:55:66"),
				createInterface(t, "wlp3s0", "aa:bb:cc:dd:ee:ff"),
			},
			expectedResult: []net.HardwareAddr{
				parseMACAddress(t, "11:22:33:44:55:66"),
				parseMACAddress(t, "aa:bb:cc:dd:ee:ff"),
			},
		},
		{
			testName:    "interface access error",
			mockError:   errAccessDenied,
			expectError: true,
		},
		{
			testName:       "empty interface list",
			mockInterfaces: []*wifi.Interface{},
			expectedResult: []net.HardwareAddr{},
		},
		{
			testName: "interface with nil address",
			mockInterfaces: []*wifi.Interface{
				createInterface(t, "wlan0", ""),
				createInterface(t, "wlan1", "11:22:33:44:55:66"),
			},
			expectedResult: []net.HardwareAddr{
				nil,
				parseMACAddress(t, "11:22:33:44:55:66"),
			},
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			t.Parallel()

			mockHandler := NewWiFiHandle(t)
			serviceInstance := wifiService.New(mockHandler)

			mockHandler.On("Interfaces").
				Return(scenario.mockInterfaces, scenario.mockError)

			addresses, err := serviceInstance.GetAddresses()

			if scenario.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
				assert.Nil(t, addresses)
			} else {
				assert.NoError(t, err)

				if expectedAddrs, ok := scenario.expectedResult.([]net.HardwareAddr); ok {
					assert.Equal(t, expectedAddrs, addresses)
				}
			}
		})
	}
}

func TestWiFiService_GetInterfaceNames(t *testing.T) {
	t.Parallel()

	testScenarios := []wifiTestCase{
		{
			testName: "successful name retrieval",
			mockInterfaces: []*wifi.Interface{
				createInterface(t, "eth0", "11:22:33:44:55:66"),
				createInterface(t, "eth1", "aa:bb:cc:dd:ee:ff"),
			},
			expectedResult: []string{"eth0", "eth1"},
		},
		{
			testName:    "driver error",
			mockError:   errHardwareMissing,
			expectError: true,
		},
		{
			testName:       "no interfaces found",
			mockInterfaces: []*wifi.Interface{},
			expectedResult: []string{},
		},
		{
			testName: "multiple interfaces",
			mockInterfaces: []*wifi.Interface{
				createInterface(t, "wlan0", "11:22:33:44:55:66"),
				createInterface(t, "wlp2s0", "aa:bb:cc:dd:ee:ff"),
				createInterface(t, "wlp3s0", "00:11:22:33:44:55"),
			},
			expectedResult: []string{"wlan0", "wlp2s0", "wlp3s0"},
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			t.Parallel()

			mockHandler := NewWiFiHandle(t)
			serviceInstance := wifiService.New(mockHandler)

			mockHandler.On("Interfaces").
				Return(scenario.mockInterfaces, scenario.mockError)

			names, err := serviceInstance.GetNames()

			if scenario.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
				assert.Nil(t, names)
			} else {
				assert.NoError(t, err)

				if expectedNames, ok := scenario.expectedResult.([]string); ok {
					assert.Equal(t, expectedNames, names)
				}
			}
		})
	}
}

func TestWiFiService_New(t *testing.T) {
	t.Parallel()

	mockHandler := NewWiFiHandle(t)
	serviceInstance := wifiService.New(mockHandler)

	assert.NotNil(t, serviceInstance)
	assert.Equal(t, mockHandler, serviceInstance.WiFi)
}

func createInterface(t *testing.T, name string, mac string) *wifi.Interface {
	t.Helper()
	var hwAddr net.HardwareAddr
	if mac != "" {
		hwAddr = parseMACAddress(t, mac)
	}

	return &wifi.Interface{
		Name:         name,
		HardwareAddr: hwAddr,
	}
}

func parseMACAddress(t *testing.T, addr string) net.HardwareAddr {
	t.Helper()
	hw, err := net.ParseMAC(addr)
	if err != nil {
		panic("invalid MAC address in test: " + addr)
	}

	return hw
}

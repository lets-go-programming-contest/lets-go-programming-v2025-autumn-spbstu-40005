package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "polina.gavrilova/task-6/internal/wifi"
)

//go:generate mockery --name=WiFiHandle --output=. --outpkg=wifi_test --case=underscore

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrDriverNotLoaded  = errors.New("driver not loaded")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		returnErr    error
		returnIfaces []*wifi.Interface
		wantErr      bool
		wantAddrs    []net.HardwareAddr
	}{
		{
			name: "success",
			returnIfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: mustParseMAC(t, "01:23:45:67:89:00")},
				{Name: "eth0", HardwareAddr: mustParseMAC(t, "ab:cd:ef:01:23:45")},
			},
			wantAddrs: []net.HardwareAddr{
				mustParseMAC(t, "01:23:45:67:89:00"),
				mustParseMAC(t, "ab:cd:ef:01:23:45"),
			},
		},
		{
			name:      "error from Interfaces",
			returnErr: ErrPermissionDenied,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &WiFiHandle{}
			service := myWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(tt.returnIfaces, tt.returnErr)

			addrs, err := service.GetAddresses()

			if tt.wantErr {
				require.ErrorContains(t, err, "getting interfaces:")
				require.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantAddrs, addrs)
			}
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		returnErr    error
		returnIfaces []*wifi.Interface
		wantErr      bool
		wantNames    []string
	}{
		{
			name: "success",
			returnIfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: mustParseMAC(t, "ab:cd:ef:01:23:45")},
				{Name: "eth1", HardwareAddr: mustParseMAC(t, "01:23:45:67:89:00")},
			},
			wantNames: []string{"wlan0", "eth1"},
		},
		{
			name:      "error from Interfaces",
			returnErr: ErrDriverNotLoaded,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &WiFiHandle{}
			service := myWifi.New(mockWiFi)

			mockWiFi.On("Interfaces").Return(tt.returnIfaces, tt.returnErr)

			names, err := service.GetNames()

			if tt.wantErr {
				require.ErrorContains(t, err, "getting interfaces:")
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}
		})
	}
}

func mustParseMAC(t *testing.T, s string) net.HardwareAddr {
	t.Helper()

	hw, err := net.ParseMAC(s)
	if err != nil {
		t.Fatalf("invalid MAC in test: %s", s)
	}

	return hw
}

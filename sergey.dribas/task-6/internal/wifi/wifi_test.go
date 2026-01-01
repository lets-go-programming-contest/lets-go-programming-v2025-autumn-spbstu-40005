// wifi/wifi_test.go
package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

// MockWiFi implements WiFiHandle
type MockWiFi struct {
	InterfacesFunc func() ([]*wifi.Interface, error)
}

func (m *MockWiFi) Interfaces() ([]*wifi.Interface, error) {
	return m.InterfacesFunc()
}

func mustParseMAC(s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return addr
}

func TestWiFiService_GetAddresses(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func() ([]*wifi.Interface, error)
		expectAddrs []net.HardwareAddr
		expectErr   error
	}{
		{
			name: "success with interfaces",
			mockFunc: func() ([]*wifi.Interface, error) {
				return []*wifi.Interface{
					{HardwareAddr: mustParseMAC("00:11:22:33:44:55")},
					{HardwareAddr: mustParseMAC("aa:bb:cc:dd:ee:ff")},
				}, nil
			},
			expectAddrs: []net.HardwareAddr{
				mustParseMAC("00:11:22:33:44:55"),
				mustParseMAC("aa:bb:cc:dd:ee:ff"),
			},
			expectErr: nil,
		},
		{
			name: "empty interfaces",
			mockFunc: func() ([]*wifi.Interface, error) {
				return []*wifi.Interface{}, nil
			},
			expectAddrs: []net.HardwareAddr{},
			expectErr:   nil,
		},
		{
			name: "error from Interfaces()",
			mockFunc: func() ([]*wifi.Interface, error) {
				return nil, errors.New("wifi error")
			},
			expectAddrs: nil,
			expectErr:   errors.New("getting interfaces: wifi error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockWiFi{InterfacesFunc: tt.mockFunc}
			service := New(mock)

			addrs, err := service.GetAddresses()

			if tt.expectErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectErr.Error())
				require.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectAddrs, addrs)
			}
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func() ([]*wifi.Interface, error)
		expectNames []string
		expectErr   error
	}{
		{
			name: "success with names",
			mockFunc: func() ([]*wifi.Interface, error) {
				return []*wifi.Interface{
					{Name: "wlan0"},
					{Name: "eth1"},
				}, nil
			},
			expectNames: []string{"wlan0", "eth1"},
			expectErr:   nil,
		},
		{
			name: "empty interfaces",
			mockFunc: func() ([]*wifi.Interface, error) {
				return []*wifi.Interface{}, nil
			},
			expectNames: []string{},
			expectErr:   nil,
		},
		{
			name: "error from Interfaces()",
			mockFunc: func() ([]*wifi.Interface, error) {
				return nil, errors.New("wifi error")
			},
			expectNames: nil,
			expectErr:   errors.New("getting interfaces: wifi error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockWiFi{InterfacesFunc: tt.mockFunc}
			service := New(mock)

			names, err := service.GetNames()

			if tt.expectErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectErr.Error())
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectNames, names)
			}
		})
	}
}

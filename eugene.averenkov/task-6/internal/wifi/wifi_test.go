package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{}
	service := New(mockHandle)

	require.Equal(t, mockHandle, service.WiFi, "Expected WiFi to be set")
}

func TestWiFiService_GetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				createMockInterface("wlan0", "00:11:22:33:44:55"),
				createMockInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
				createMockInterface("eth0", "11:22:33:44:55:66"),
			}, nil
		},
	}

	service := New(mockHandle)
	addresses, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, addresses, 3, "Expected 3 addresses")

	expected1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	expected2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)
	expected3, err := net.ParseMAC("11:22:33:44:55:66")
	require.NoError(t, err)

	require.Equal(t, expected1.String(), addresses[0].String(), "First address mismatch")
	require.Equal(t, expected2.String(), addresses[1].String(), "Second address mismatch")
	require.Equal(t, expected3.String(), addresses[2].String(), "Third address mismatch")
}

func TestWiFiService_GetAddresses_Empty(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{}, nil
		},
	}

	service := New(mockHandle)
	addresses, err := service.GetAddresses()
	require.NoError(t, err)
	require.Empty(t, addresses, "Expected empty slice")
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errors.New("interface error")
		},
	}

	service := New(mockHandle)
	addresses, err := service.GetAddresses()
	require.Error(t, err, "Expected error")
	require.Nil(t, addresses, "Expected nil result on error")
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				createMockInterface("wlan0", "00:11:22:33:44:55"),
				createMockInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
				createMockInterface("eth0", "11:22:33:44:55:66"),
				createMockInterface("wlp3s0", "22:33:44:55:66:77"),
			}, nil
		},
	}

	service := New(mockHandle)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Len(t, names, 4, "Expected 4 names")

	expected := []string{"wlan0", "wlan1", "eth0", "wlp3s0"}
	for i, name := range names {
		require.Equal(t, expected[i], name, "Name mismatch at index %d", i)
	}
}

func TestWiFiService_GetNames_Empty(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{}, nil
		},
	}

	service := New(mockHandle)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names, "Expected empty slice")
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errors.New("interface error")
		},
	}

	service := New(mockHandle)
	names, err := service.GetNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
}

func TestWiFiService_BothMethodsSameData(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				createMockInterface("wlan0", "00:11:22:33:44:55"),
				createMockInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
				createMockInterface("eth0", "11:22:33:44:55:66"),
			}, nil
		},
	}

	service := New(mockHandle)

	addresses, err := service.GetAddresses()
	require.NoError(t, err, "GetAddresses error")
	require.Len(t, addresses, 3, "Expected 3 addresses")

	names, err := service.GetNames()
	require.NoError(t, err, "GetNames error")
	require.Len(t, names, 3, "Expected 3 names")

	require.Equal(t, "wlan0", names[0], "First interface name mismatch")
	require.Equal(t, "wlan1", names[1], "Second interface name mismatch")
	require.Equal(t, "eth0", names[2], "Third interface name mismatch")
}

func TestWiFiService_NilHardwareAddr(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{
					Name:         "wlan0",
					HardwareAddr: nil,
				},
				createMockInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
			}, nil
		},
	}

	service := New(mockHandle)
	addresses, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, addresses, 2, "Expected 2 addresses")
	require.Nil(t, addresses[0], "Expected first address to be nil")

	expected, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)
	require.Equal(t, expected.String(), addresses[1].String(), "Second address mismatch")
}

func TestWiFiService_InterfaceWithEmptyName(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{
					Name:         "",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				},
				createMockInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
			}, nil
		},
	}

	service := New(mockHandle)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Len(t, names, 2, "Expected 2 names")
	require.Equal(t, "", names[0], "Expected empty name")
	require.Equal(t, "wlan1", names[1], "Expected wlan1 name")
}

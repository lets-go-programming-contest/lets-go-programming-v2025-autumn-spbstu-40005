package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

type mockWiFiHandle struct {
	interfacesFunc func() ([]*wifi.Interface, error)
}

func (m *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	if m.interfacesFunc != nil {
		return m.interfacesFunc()
	}
	return nil, errors.New("not implemented")
}

func createMockInterface(name string, addr string) *wifi.Interface {
	hwAddr, _ := net.ParseMAC(addr)
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: hwAddr,
	}
}

func TestWiFiHandle_Mock_Interfaces_Success(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				createMockInterface("wlan0", "00:11:22:33:44:55"),
				createMockInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
			}, nil
		},
	}

	interfaces, err := mockHandle.Interfaces()
	require.NoError(t, err)
	require.Len(t, interfaces, 2)
	require.Equal(t, "wlan0", interfaces[0].Name)
	require.Equal(t, "wlan1", interfaces[1].Name)
}

func TestWiFiHandle_Mock_Interfaces_Error(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return nil, errors.New("mock error")
		},
	}

	interfaces, err := mockHandle.Interfaces()
	require.Error(t, err)
	require.Nil(t, interfaces)
}

func TestWiFiHandle_Mock_Interfaces_Empty(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{}, nil
		},
	}

	interfaces, err := mockHandle.Interfaces()
	require.NoError(t, err)
	require.Empty(t, interfaces)
}

func TestWiFiHandle_Mock_Interfaces_Default(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{}
	interfaces, err := mockHandle.Interfaces()
	require.Error(t, err)
	require.Nil(t, interfaces)
}

func TestWiFiHandle_CreateMockInterface(t *testing.T) {
	t.Parallel()

	iface := createMockInterface("eth0", "11:22:33:44:55:66")
	require.Equal(t, "eth0", iface.Name)

	expectedMAC, err := net.ParseMAC("11:22:33:44:55:66")
	require.NoError(t, err)
	require.Equal(t, expectedMAC.String(), iface.HardwareAddr.String())
}

func TestWiFiHandle_Mock_InterfaceWithNilMAC(t *testing.T) {
	t.Parallel()

	mockHandle := &mockWiFiHandle{
		interfacesFunc: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{
					Name:         "eth0",
					HardwareAddr: nil,
				},
			}, nil
		},
	}

	interfaces, err := mockHandle.Interfaces()
	require.NoError(t, err)
	require.Len(t, interfaces, 1)
	require.Equal(t, "eth0", interfaces[0].Name)
	require.Nil(t, interfaces[0].HardwareAddr)
}

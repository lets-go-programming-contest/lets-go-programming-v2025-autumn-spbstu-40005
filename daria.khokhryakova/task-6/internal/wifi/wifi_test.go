package wifi_test

import (
	"net"
	"testing"

	myWifi "github.com/DariaKhokhryakova/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type wifiHandleMock struct {
	mock.Mock
}

func (m *wifiHandleMock) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var result []*wifi.Interface
	if args.Get(0) != nil {
		result = args.Get(0).([]*wifi.Interface)
	}

	return result, args.Error(1)
}

func createTestInterface(name string, addr string) *wifi.Interface {
	hwAddr, _ := net.ParseMAC(addr)
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: hwAddr,
	}
}

func TestWiFiService_GetAddresses_Success(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	interfaces := []*wifi.Interface{
		createTestInterface("wlan0", "00:11:22:33:44:55"),
		createTestInterface("wlan1", "aa:bb:cc:dd:ee:ff"),
	}
	mock.On("Interfaces").Return(interfaces, nil)

	service := myWifi.New(mock)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 2)
	assert.Equal(t, "00:11:22:33:44:55", addrs[0].String())
	assert.Equal(t, "aa:bb:cc:dd:ee:ff", addrs[1].String())
	mock.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_InterfaceFetchError(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	mock.On("Interfaces").Return(([]*wifi.Interface)(nil), assert.AnError)

	service := myWifi.New(mock)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting interfaces:")
	assert.Nil(t, addrs)
	mock.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_EmptyList(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	mock.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := myWifi.New(mock)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Empty(t, addrs)
	mock.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_NullMACAddress(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	interfaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: nil,
		},
	}
	mock.On("Interfaces").Return(interfaces, nil)

	service := myWifi.New(mock)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 1)
	assert.Nil(t, addrs[0])
	mock.AssertExpectations(t)
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	interfaces := []*wifi.Interface{
		createTestInterface("wlan0", "00:11:22:33:44:55"),
		createTestInterface("eth1", "11:22:33:44:55:66"),
		createTestInterface("wifi2", "22:33:44:55:66:77"),
	}
	mock.On("Interfaces").Return(interfaces, nil)

	service := myWifi.New(mock)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"wlan0", "eth1", "wifi2"}, names)
	mock.AssertExpectations(t)
}

func TestWiFiService_GetNames_InterfaceFetchError(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	mock.On("Interfaces").Return(([]*wifi.Interface)(nil), assert.AnError)

	service := myWifi.New(mock)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting interfaces:")
	assert.Nil(t, names)
	mock.AssertExpectations(t)
}

func TestWiFiService_GetNames_EmptyList(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	mock.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := myWifi.New(mock)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Empty(t, names)
	mock.AssertExpectations(t)
}

func TestWiFiService_GetNames_BlankInterfaceName(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	interfaces := []*wifi.Interface{
		{
			Name:         "",
			HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		},
	}
	mock.On("Interfaces").Return(interfaces, nil)

	service := myWifi.New(mock)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{""}, names)
	mock.AssertExpectations(t)
}

func TestWiFiService_Initialization(t *testing.T) {
	t.Parallel()

	mock := &wifiHandleMock{}
	service := myWifi.New(mock)

	assert.NotNil(t, service)
	assert.Equal(t, mock, service.WiFi)
}

// +build hal_atlas300 hal_atlas500

package acldevice

import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor"
	_ "gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor/acldevice/include"
)

type acldeviceMonitor struct {
}

func (a *acldeviceMonitor) Init() error {
	return nil
}

func (a *acldeviceMonitor) BindDevice(id int) error {
	return nil
}

func (a *acldeviceMonitor) GetDeviceCount() (int, error) {
	return 1, nil
}

// GetMemInfo return ascend mem info in bytes
func (a *acldeviceMonitor) GetMemInfo(id int) (free int, total int, err error) {
	//ascend dsmi return mem info in MB
	use, total, err := getMemoryInfo(id)
	if err != nil {
		return 0, 0, err
	}
	return (total - use) * monitor.OneMB, total * monitor.OneMB, nil
}

func (a *acldeviceMonitor) GetSDKProperties() (*monitor.SDKProperties, error) {
	return &monitor.SDKProperties{
		SDKMajor: 0,
		SDKMinor: 0,
	}, nil
}

func (a *acldeviceMonitor) GetDeviceProperties(id int) (*monitor.DeviceProperties, error) {
	name, err := getDeviceName(id)
	if err != nil {
		return nil, err
	}
	return &monitor.DeviceProperties{
		HardwareMajor: 0,
		HardwareMinor: 0,
		DeviceName:    name,
	}, nil
}

func (a *acldeviceMonitor) UnbindDevice(id int) {
}

func (a *acldeviceMonitor) UnInit() error {
	return nil
}

func init() {
	am := &acldeviceMonitor{}
	if err := am.Init(); err != nil {
		panic(err)
	}
	monitor.DeviceMonitors[monitor.ATLAS] = am
	monitor.DefaultMonitor = am
	monitor.DefaultDevice = monitor.ATLAS
}

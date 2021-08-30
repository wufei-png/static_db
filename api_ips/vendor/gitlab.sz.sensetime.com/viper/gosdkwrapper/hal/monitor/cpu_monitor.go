package monitor

import (
	"runtime"
)

type CPUMonitor struct {
}

func (c *CPUMonitor) Init() error {
	return nil
}

func (c *CPUMonitor) GetDeviceCount() (int, error) {
	return 1, nil
}

//TODO use runtime.MemStats to get meminfo?
func (c *CPUMonitor) GetMemInfo(id int) (free int, total int, err error) {
	return 100000, 100000, nil
}

func (c *CPUMonitor) GetSDKProperties() (*SDKProperties, error) {
	return &SDKProperties{
		SDKMajor: 0,
		SDKMinor: 0,
	}, nil
}

func (c *CPUMonitor) GetDeviceProperties(id int) (*DeviceProperties, error) {
	return &DeviceProperties{
		HardwareMajor: 0,
		HardwareMinor: 0,
		DeviceName:    runtime.GOARCH,
	}, nil
}

func (c *CPUMonitor) BindDevice(id int) error {
	return nil
}

func (c *CPUMonitor) UnbindDevice(id int) {

}

func (c *CPUMonitor) UnInit() error {
	return nil
}

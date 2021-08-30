// +build hal_stpu

package stpu

import (
	"sync"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor"
)

type stpuMonitor struct {
	once sync.Once
}

func (m *stpuMonitor) Init() error {
	return nil
}

func (m *stpuMonitor) BindDevice(id int) error {
	if err := m.initHalTop(id); err != nil {
		return err
	}
	if err := m.openHalDevice(id); err != nil {
		return err
	}
	m.once.Do(m.initHalLibML)

	return nil
}

func (m *stpuMonitor) GetDeviceCount() (int, error) {
	return 1, nil
}

// GetMemInfo return stpu mem info in bytes
func (m *stpuMonitor) GetMemInfo(id int) (free int, total int, err error) {

	//stpuhal return stpu meminfo in MB
	use, total, err := m.getPuMem()
	if err != nil {
		return 0, 0, err
	}
	return (total - use) * monitor.OneMB, total * monitor.OneMB, nil
}

func (m *stpuMonitor) GetSDKProperties() (*monitor.SDKProperties, error) {
	return &monitor.SDKProperties{
		SDKMajor: 0,
		SDKMinor: 0,
	}, nil
}

func (m *stpuMonitor) GetDeviceProperties(id int) (*monitor.DeviceProperties, error) {
	deviceName, err := m.getProductModelNumber()
	if err != nil {
		return nil, err
	}
	return &monitor.DeviceProperties{
		HardwareMajor: 0,
		HardwareMinor: 0,
		DeviceName:    deviceName,
	}, nil
}

func (m *stpuMonitor) UnbindDevice(id int) {
	m.closeHalDevice(id)
	m.deinitHalTop()
}

func (m *stpuMonitor) UnInit() error {
	return nil
}

func init() {
	//sync.Once is for initHalLibML which only init once
	sm := &stpuMonitor{once: sync.Once{}}
	if err := sm.Init(); err != nil {
		panic(err)
	}
	monitor.DeviceMonitors[monitor.STPU] = sm
	monitor.DefaultMonitor = sm
	monitor.DefaultDevice = monitor.STPU
}

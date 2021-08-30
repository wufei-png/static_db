// +build hal_cambricon

package cambricon

import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor"
)

func init() {
	cm := &cambriconDeviceMonitorImpl{}
	if err := cm.Init(); err != nil {
		panic(err)
	}
	monitor.DeviceMonitors[monitor.CAMBRICON] = cm
	monitor.DefaultMonitor = cm
	monitor.DefaultDevice = monitor.CAMBRICON
}

type cambriconDeviceMonitorImpl struct{}

func (c *cambriconDeviceMonitorImpl) Init() error {
	return cnInit()
}

func (c *cambriconDeviceMonitorImpl) BindDevice(_ int) error {
	return nil
}

func (c *cambriconDeviceMonitorImpl) UnbindDevice(deviceId int) {
}

func (c *cambriconDeviceMonitorImpl) GetDeviceCount() (int, error) {
	return cnrtGetDeviceCount()
}

func (c *cambriconDeviceMonitorImpl) GetSDKProperties() (*monitor.SDKProperties, error) {
	ret := &monitor.SDKProperties{}
	ret.SDKMajor, ret.SDKMinor = cnrtGetLibVersion()

	return ret, nil
}

// GetMemInfo return mlu mem info in bytes
func (c *cambriconDeviceMonitorImpl) GetMemInfo(deviceId int) (free int, total int, err error) {
	use, total, err := cndevGetMemoryUsage(deviceId)
	if err != nil {
		return 0, 0, err
	}
	return (total - use) * monitor.OneMB, total * monitor.OneMB, nil
}

func (c *cambriconDeviceMonitorImpl) GetDeviceProperties(_ int) (*monitor.DeviceProperties, error) {
	deviceName := cndevCardName()
	return &monitor.DeviceProperties{DeviceName: deviceName}, nil
}

func (c *cambriconDeviceMonitorImpl) UnInit() error {
	return cnUnInit()
}
